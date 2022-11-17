package filter

import (
	"fmt"
	"strconv"
	"strings"
)

// FilterExpr represents the top level expression.
type FilterExpr struct {
	expr Expr
}

// String returns an indented, pretty-printed version of the parsed
// program.
func (p *FilterExpr) String() string {
	return fmt.Sprintf("%s\n\n", p.expr.String())
}

// Expr is the abstract syntax tree for any expression.
type Expr interface {
	String() string
}

// binaryExpr is an expression like 1 + 2.
type binaryExpr struct {
	Left  Expr
	Op    Token
	Right Expr
}

func (e *binaryExpr) String() string {
	var opStr string
	opStr = " " + e.Op.String() + " "
	return "(" + e.Left.String() + opStr + e.Right.String() + ")"
}

// strExpr is a literal string like "foo".
type strExpr struct {
	Value string
}

func (e *strExpr) String() string {
	return strconv.Quote(e.Value)
}

// varExpr is a variable reference (name, description,location).
type varExpr struct {
	Name string
}

func (v *varExpr) String() string {
	return strconv.Quote(v.Name)
}

type listExpr struct {
	Items []string
}

func (v *listExpr) String() string {
	return fmt.Sprintf("[%s]", strings.Join(v.Items, ","))
}
