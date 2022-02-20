package search

import (
	"fmt"
	"regexp"
	"time"
)

// Grammar
//
// expression: logical | (logical)												;
// logical: equality | equality (( "&" | "|" ) equality)*						;
// equality: term ( ("==" | "!=" | "<" | "<=" | ">" | ">=" | "~") primary )*	;
// term: VAR_NAME																;
// primary: STRING | DATE | REGEX												;
//
// Note: to make it easy for user to enter expressions "&" and "|" are equivalent with "&&" and "||"
// Date has only one format accepted: 01/02/2002 ( 02 -> month )
// Regex format is /regex/ and it has to be Go regex.

// ParseError (actually *ParseError) is the type of error returned by ParseSearchExpression.
type ParseError struct {
	// Source line/column position where the error occurred.
	Position int
	// Error message.
	Message string
}

// Error returns a formatted version of the error, including the line
// and column numbers.
func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at %d: %s", e.Position, e.Message)
}

type parser struct {
	// Lexer instance and current token values
	lexer *Lexer
	pos   int    // position of last token (tok)
	tok   Token  // last lexed token
	val   string // string value of last token (or "")
}

func parse(src []byte) (searchExpr *BinaryExpr, err error) {
	defer func() {
		// The parser uses panic with a *ParseError to signal parsing
		// errors internally, and they're caught here. This
		// significantly simplifies the recursive descent calls as
		// we don't have to check errors everywhere.
		if r := recover(); r != nil {
			// Convert to ParseError or re-panic
			err = r.(*ParseError)
		}
	}()
	lexer := NewLexer(src)
	p := parser{lexer: lexer}
	p.next() // initialize p.tok

	// Parse into abstract syntax tree
	searchExpr = p.expression().(*BinaryExpr)

	return
}

// Parse a logic expression
//
// equality | equality (( "&" | "|" ) equality)*
//
func (p *parser) expression() Expr {
	var expr Expr

	if p.matches(LPAREN) {
		p.next()
		expr = p.expression()
		p.consume(RPAREN, "expected ')' after expression")
	} else {
		expr = p.equality()
	}

	for p.matches(AND, OR) {
		op := p.tok
		p.next()

		var right Expr
		if p.matches(LPAREN) {
			p.next()
			right = p.expression()
			p.consume(RPAREN, "expected ')' after expression")
		} else {
			right = p.equality()
		}

		expr = &BinaryExpr{Left: expr, Op: op, Right: right}
	}

	return expr
}

// Parse equality expression
//
// term ( ("==" | "!=" | "<" | "<=" | ">" | ">=" | "~") primary )*	;
//
func (p *parser) equality() Expr {
	p.expect(VAR_NAME)
	expr := &BinaryExpr{Left: p.primary()}

	switch p.tok {
	case GREATER, GTE, LESS, LTE, EQUALS, NOT_EQUALS, TILDA:
		expr.Op = p.tok
		p.next()
	default:
		panic(p.errorf("expected operator instead of %s", p.tok))

	}

	expr.Right = p.primary()

	return expr
}

// Parse primary
//
// STRING | DATE | REGEX
//
func (p *parser) primary() Expr {
	var expr Expr

	switch p.tok {
	case VAR_NAME:
		expr = &VarExpr{p.val}
	case STRING:
		expr = &StrExpr{p.val}
	case DATE:
		expr = p.dateExpr(p.val)
	case REGEX:
		expr = p.regexExpr(p.val)
	default:
		panic(p.errorf("expected string instead of %s", p.tok))
	}

	p.next()

	return expr
}

func (p *parser) dateExpr(date string) Expr {
	t, err := time.Parse("02/01/2006", date)
	if err != nil {
		panic(p.errorf("expected date instead of '%s'", date))
	}

	return &DateExpr{t}
}

func (p *parser) regexExpr(regex string) Expr {
	r, err := regexp.Compile(regex)
	if err != nil {
		panic(p.errorf("expected valid regex instead of '%s'", regex))
	}

	return &RegexExpr{r}
}

// Parse next token into p.tok (and set p.pos and p.val).
func (p *parser) next() {
	p.pos, p.tok, p.val = p.lexer.Scan()
	if p.tok == ILLEGAL {
		panic(p.errorf("%s", p.val))
	}
}

// Return true iff current token matches one of the given operators,
// but don't parse next token.
func (p *parser) matches(operators ...Token) bool {
	for _, operator := range operators {
		if p.tok == operator {
			return true
		}
	}
	return false
}

// Ensure current token is tok, and parse next token into p.tok.
func (p *parser) expect(tok Token) {
	if p.tok != tok {
		panic(p.errorf("expected %s instead of %s", tok, p.tok))
	}
}

func (p *parser) consume(tok Token, msg string) {
	if !p.matches(tok) {
		panic(p.errorf(msg))
	}
	p.next()
}

// Format given string and args with Sprintf and return an error
// with that message and the current position.
func (p *parser) errorf(format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	return &ParseError{p.pos, message}
}
