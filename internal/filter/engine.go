package filter

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/tupyy/gophoto/internal/entity"
)

var (
	// WrongOpError means that the wrong operator is used in comparison.
	WrongOpError = errors.New("wrong op")
	// FieldNotFoundError means the keyword used in filter expression is not a field of the album.
	FieldNotFoundError = errors.New("album field not found")
)

type Filter struct {
	expr *binaryExpr
}

func New(filterExpr string) (*Filter, error) {
	expr, err := parse([]byte(filterExpr))
	if err != nil {
		return nil, err
	}

	return &Filter{expr}, nil
}

// Resolve tries to resolve the album against the filter expression.
// Returns false if the album does not pass the expression.
func (f *Filter) Resolve(album entity.Album) (bool, error) {
	return resolveAST(f.expr, album)
}

func resolveAST(rootExpr *binaryExpr, album entity.Album) (bool, error) {
	var (
		leftResult  bool
		rightResult bool
		err         error
	)

	if left, ok := rootExpr.Left.(*binaryExpr); ok {
		leftResult, err = resolveAST(left, album)
	} else {
		leftResult, err = resolveExpr(rootExpr, album)
	}

	if err != nil {
		return false, err
	}

	if right, ok := rootExpr.Right.(*binaryExpr); ok {
		rightResult, err = resolveAST(right, album)
	} else {
		rightResult, err = resolveExpr(rootExpr, album)
	}

	if err != nil {
		return false, err
	}

	var result bool
	switch rootExpr.Op {
	case AND:
		result = leftResult && rightResult
	case OR:
		result = leftResult || rightResult
	default:
		return leftResult, err // leftResult is the same as rightResult
	}

	return result, nil

}

func resolveExpr(expr *binaryExpr, album entity.Album) (bool, error) {
	variable := expr.Left.(*varExpr)

	// if it is a date, try to parse the value
	switch variable.Name {
	case "date":
		return resolveDate(expr, album)
	case "tag":
		return resolveArrayField(expr, album, func(album entity.Album) []string {
			list := make([]string, 0, len(album.Tags))
			for _, t := range album.Tags {
				list = append(list, t.Name)
			}
			return list
		})
	case "permissions.user":
		return resolveArrayField(expr, album, func(album entity.Album) []string {
			list := make([]string, 0, len(album.UserPermissions))
			for _, u := range album.UserPermissions {
				list = append(list, u.OwnerID)
			}
			return list
		})
	case "permissions.group":
		return resolveArrayField(expr, album, func(album entity.Album) []string {
			list := make([]string, 0, len(album.GroupPermissions))
			for _, g := range album.GroupPermissions {
				list = append(list, g.OwnerID)
			}
			return list
		})
	default:
		return resolveCommonField(expr, album)
	}
}

func resolveDate(expr *binaryExpr, album entity.Album) (bool, error) {
	dateExpr, ok := expr.Right.(*strExpr)
	if !ok {
		return false, fmt.Errorf("expect string got '%s'", expr.Right.String())
	}

	date, err := time.Parse("02/01/2006", dateExpr.Value)
	if err != nil {
		return false, fmt.Errorf("expected date instead of '%s'", dateExpr.Value)
	}

	switch expr.Op {
	case GREATER, GTE:
		return album.CreatedAt.After(date), nil
	case LTE, LESS:
		return album.CreatedAt.Before(date), nil
	case EQUALS:
		return date.Day() == album.CreatedAt.Day() && date.Month() == album.CreatedAt.Month() && date.Year() == album.CreatedAt.Year(), nil
	case NOT_EQUALS:
		return date.Day() != album.CreatedAt.Day() || date.Month() != album.CreatedAt.Month() || date.Year() != album.CreatedAt.Year(), nil
	default:
		return false, WrongOpError
	}
}

func resolveArrayField(expr *binaryExpr, album entity.Album, getItemsFn func(album entity.Album) []string) (bool, error) {
	value, ok := expr.Right.(*strExpr)
	if !ok {
		return false, fmt.Errorf("expect string got '%s'", expr.Right.String())
	}

	switch expr.Op {
	case EQUALS:
		return strings.Index(strings.Join(getItemsFn(album), " "), value.Value) >= 0, nil
	case NOT_EQUALS:
		return strings.Index(strings.Join(getItemsFn(album), " "), value.Value) == -1, nil
	}

	return false, fmt.Errorf("%w tag comparison cannot have something else than '=' or '!='.got '%s'", WrongOpError, expr.Op)
}

func resolveCommonField(expr *binaryExpr, album entity.Album) (bool, error) {
	var varValue string
	variable := expr.Left.(*varExpr)

	switch variable.Name {
	case "name":
		varValue = album.Name
	case "description":
		varValue = album.Description
	case "location":
		varValue = album.Location
	case "owner":
		varValue = album.Owner
	default:
		return false, fmt.Errorf("%w unknown field %s", FieldNotFoundError, variable.Name)
	}

	if expr.Op == IN {
		listExr, ok := expr.Right.(*listExpr)
		if !ok {
			return false, fmt.Errorf("expect list got '%s'", expr.Right.String())
		}
		for _, item := range listExr.Items {
			if item == varValue {
				return true, nil
			}
		}
		return false, nil
	}

	// we expect a string here.
	value, ok := expr.Right.(*strExpr)
	if !ok {
		return false, fmt.Errorf("expect string got '%s'", expr.Right.String())
	}

	switch expr.Op {
	case EQUALS:
		return varValue == value.Value, nil
	case NOT_EQUALS:
		return varValue != value.Value, nil
	case GREATER:
		return varValue > value.Value, nil
	case GTE:
		return varValue >= value.Value, nil
	case LESS:
		return varValue < value.Value, nil
	case LTE:
		return varValue <= value.Value, nil
	case LIKE:
		regex, ok := expr.Right.(*strExpr)
		if !ok {
			return false, fmt.Errorf("expected string got '%s'", expr.Right.String())
		}
		rxp, err := regexp.Compile(regex.Value)
		if err != nil {
			return false, fmt.Errorf("failed to compile pattern '%s': %s", regex.Value, err)
		}
		return rxp.MatchString(varValue), nil
	default:
		return false, fmt.Errorf("%w unaccepted operator used in common fields comparison. got '%s'", WrongOpError, expr.Op)
	}
}
