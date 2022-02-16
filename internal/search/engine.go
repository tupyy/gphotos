package search

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tupyy/gophoto/internal/entity"
)

var (
	// WrongOpError means that the wrong operator is used in comparison.
	WrongOpError = errors.New("wrong op")
	// FieldNotFoundError means the keyword used in filter expression is not a field of the album.
	FieldNotFoundError = errors.New("album field not found")
	// DateParseError means the value of the date variable cannot be parsed.
	DateParseError = errors.New("date parse error")

	// accepted date time layouts
	dateLayouts = []string{"02/Jan/2006", "02/01/2006"}
)

type SearchEngine struct {
	expr *BinaryExpr
}

func NewSearchEngine(filterExpr string) (*SearchEngine, error) {
	expr, err := parseSearchExpression([]byte(filterExpr))
	if err != nil {
		return nil, err
	}

	return &SearchEngine{expr}, nil
}

// Resolve tries to resolve the album against the search expression.
// Returns false if the album does not pass the expression.
func (f *SearchEngine) Resolve(album entity.Album) (bool, error) {
	return resolveAST(f.expr, album)
}

func resolveAST(rootExpr *BinaryExpr, album entity.Album) (bool, error) {
	var (
		leftResult  bool
		rightResult bool
		err         error
	)

	left, hasLeft := rootExpr.Left.(*BinaryExpr)
	if hasLeft {
		leftResult, err = resolveAST(left, album)
	} else {
		leftResult, err = resolveExpr(rootExpr, album)
	}

	if err != nil {
		return false, err
	}

	right, hasRight := rootExpr.Right.(*BinaryExpr)
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

func resolveExpr(expr *BinaryExpr, album entity.Album) (bool, error) {
	variable := expr.Left.(*VarExpr)

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

func parseTime(date string) (time.Time, error) {
	for _, layout := range dateLayouts {
		t, err := time.Parse(layout, date)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, DateParseError
}

func resolveDate(expr *BinaryExpr, album entity.Album) (bool, error) {
	value := expr.Right.(*StrExpr)

	t, err := parseTime(value.Value)
	if err != nil {
		return false, fmt.Errorf("%w accepted formats are: 02/Jan/2006 or 02/01/2006. got: %s", DateParseError, value.Value)
	}

	switch expr.Op {
	case GREATER, GTE:
		return album.CreatedAt.After(t), nil
	case LTE, LESS:
		return album.CreatedAt.Before(t), nil
	case EQUALS:
		return t.Day() == album.CreatedAt.Day() && t.Month() == album.CreatedAt.Month() && t.Year() == album.CreatedAt.Year(), nil
	case NOT_EQUALS:
		return t.Day() != album.CreatedAt.Day() || t.Month() != album.CreatedAt.Month() || t.Year() != album.CreatedAt.Year(), nil
	default:
		return false, WrongOpError
	}
}

func resolveTags(expr *BinaryExpr, album entity.Album) (bool, error) {
	value := expr.Right.(*StrExpr)

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

func resolveCommonField(expr *BinaryExpr, album entity.Album) (bool, error) {
	var varValue string
	variable := expr.Left.(*VarExpr)
	value := expr.Right.(*StrExpr)

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
