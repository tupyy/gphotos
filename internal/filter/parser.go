// Grammar
//
// expression: equality | equality ( "&&" | "||" ) equality                     ;
// equality: variable ("=" | "!=" | "<" | "<=" | ">" | ">=" | "~") primary      ;
// primary: STRING | ARRAY
//

package filter

import (
	"fmt"
)

// ParseError (actually *ParseError) is the type of error returned by parse.
type ParseError struct {
	// Source line/column position where the error occurred.
	Position int
	// Error message.
	Message string
}

// Error returns a formatted version of the error, including the line number.
func (e ParseError) Error() string {
	return fmt.Sprintf("parse error at %d: %s", e.Position, e.Message)
}

type parser struct {
	// Lexer instance and current token values
	lexer *lexer
	pos   int    // position of last token (tok)
	tok   Token  // last lexed token
	val   string // string value of last token (or "")
}

func parse(src []byte) (filterExpr *binaryExpr, err error) {
	defer func() {
		if r := recover(); r != nil {
			// Convert to ParseError or re-panic
			err = r.(ParseError)
		}
	}()

	lexer := newLexer(src)
	p := parser{lexer: lexer}
	p.next() // initialize p.tok

	// Parse into abstract syntax tree
	filterExpr = p.expression().(*binaryExpr)

	return
}

// Parse a logic expression
//
// equality | equality ( "&&" | "||" ) equality*
//
func (p *parser) expression() Expr {
	var expr Expr
	expr = p.equality()

	if !p.matches(AND, OR) && !p.matches(EOL) {
		panic(p.errorf("unexpected expression after '%s'", p.tok))
	}

	for p.matches(AND, OR) {
		op := p.tok
		p.next()
		right := p.equality()
		expr = &binaryExpr{Left: expr, Op: op, Right: right}
	}

	return expr
}

// Parse equality expression
//
// term ("==" | "!=" | "<" | "<=" | ">" | ">=" | "~") primary
//
func (p *parser) equality() Expr {
	p.expect(VARIABLE)
	expr := &binaryExpr{Left: p.primary()}

	switch p.tok {
	case GREATER, GTE, LESS, LTE, EQUALS, NOT_EQUALS, LIKE, IN:
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
	case LBRACKET:
		items := []string{}
		keepReading := true
		for keepReading {
			p.next()
			p.expect(STRING)
			items = append(items, p.val)
			p.next()
			switch p.tok {
			case RBRACKET:
				expr = &listExpr{items}
				keepReading = false
			case COMMA:
			default:
				panic(p.errorf("unexpected either string or ] instead of '%s'", p.tok))
			}
		}
	case VARIABLE:
		expr = &varExpr{p.val}
	case STRING:
		expr = &strExpr{p.val}
	default:
		panic(p.errorf("expected string instead of %s", p.tok))
	}

	p.next()

	return expr
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
	return ParseError{p.pos, message}
}
