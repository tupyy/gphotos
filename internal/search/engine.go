package search

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tupyy/gophoto/internal/entity"
)

var (
	// WrongOpError means that the wrong operator is used in comparison.
	WrongOpError = errors.New("wrong op")
	// FieldNotFoundError means the keyword used in filter expression is not a field of the album.
	FieldNotFoundError = errors.New("album field not found")
)

type FilterEngine struct {
	expr *binaryExpr
}

func NewSearchEngine(filterExpr string) (*FilterEngine, error) {
	expr, err := parse([]byte(filterExpr))
	if err != nil {
		return nil, err
	}

	return &FilterEngine{expr}, nil
}

// Resolve tries to resolve the album against the filter expression.
// Returns false if the album does not pass the expression.
func (f *FilterEngine) Resolve(album entity.Album) (bool, error) {
	return resolveAST(f.expr, album)
}

func resolveAST(rootExpr *binaryExpr, album entity.Album) (bool, error) {
	var (
		leftResult  bool
		rightResult bool
		err         error
	)

	left, hasLeft := rootExpr.Left.(*binaryExpr)
	if hasLeft {
		leftResult, err = resolveAST(left, album)
	} else {
		leftResult, err = resolveExpr(rootExpr, album)
	}

	if err != nil {
		return false, err
	}

	right, hasRight := rootExpr.Right.(*binaryExpr)
	if hasRight {
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
		return resolveTags(expr, album)
	default:
		return resolveCommonField(expr, album)
	}
}

func resolveDate(expr *binaryExpr, album entity.Album) (bool, error) {
	dateExpr, ok := expr.Right.(*dateExpr)
	if !ok {
		return false, fmt.Errorf("expect date got '%s'", expr.Right.String())
	}

	switch expr.Op {
	case GREATER, GTE:
		return album.CreatedAt.After(dateExpr.Date), nil
	case LTE, LESS:
		return album.CreatedAt.Before(dateExpr.Date), nil
	case EQUALS:
		return dateExpr.Date.Day() == album.CreatedAt.Day() && dateExpr.Date.Month() == album.CreatedAt.Month() && dateExpr.Date.Year() == album.CreatedAt.Year(), nil
	case NOT_EQUALS:
		return dateExpr.Date.Day() != album.CreatedAt.Day() || dateExpr.Date.Month() != album.CreatedAt.Month() || dateExpr.Date.Year() != album.CreatedAt.Year(), nil
	default:
		return false, WrongOpError
	}
}

func resolveTags(expr *binaryExpr, album entity.Album) (bool, error) {
	value, ok := expr.Right.(*strExpr)
	if !ok {
		return false, fmt.Errorf("expect string got '%s'", expr.Right.String())
	}

	var tags string
	for _, t := range album.Tags {
		tags = fmt.Sprintf("%s %s", tags, t.Name)
	}

	switch expr.Op {
	case EQUALS:
		return strings.Index(tags, value.Value) >= 0, nil
	case NOT_EQUALS:
		return strings.Index(tags, value.Value) == -1, nil
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
	default:
		return false, fmt.Errorf("%w unknown field %s", FieldNotFoundError, variable.Name)
	}

	if expr.Op == TILDA {
		regex, ok := expr.Right.(*regexExpr)
		if !ok {
			return false, fmt.Errorf("expected regex got '%s'", expr.Right.String())
		}

		return regex.Regex.MatchString(varValue), nil
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
	default:
		return false, fmt.Errorf("%w unaccepted operator used in common fields comparison. got '%s'", WrongOpError, expr.Op)
	}
}
