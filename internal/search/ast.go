package search

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
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

// BinaryExpr is an expression like 1 + 2.
type BinaryExpr struct {
	Left  Expr
	Op    Token
	Right Expr
}

func (e *BinaryExpr) String() string {
	var opStr string
	opStr = " " + e.Op.String() + " "
	return "(" + e.Left.String() + opStr + e.Right.String() + ")"
}

// StrExpr is a literal string like "foo".
type StrExpr struct {
	Value string
}

func (e *StrExpr) String() string {
	return strconv.Quote(e.Value)
}

// Date expression
type DateExpr struct {
	date time.Time
}

func (d *DateExpr) String() string {
	return strconv.Quote(d.date.Format("02/01/2006"))
}

type RegexExpr struct {
	regex *regexp.Regexp
}

func (r *RegexExpr) String() string {
	return strconv.Quote(r.regex.String())
}

// VarExpr is a variable reference (name, description,location).
type VarExpr struct {
	Name string
}

func (v *VarExpr) String() string {
	return strconv.Quote(v.Name)
}
