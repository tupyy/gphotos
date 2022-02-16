package search

import (
	"fmt"
)

// Parser parses logical and comparison expressions like
// `album_field_name` = 'something' & 'album_field_name2' = 'something else'
type parser struct {
	// Lexer instance and current token values
	lexer *Lexer
	pos   int    // position of last token (tok)
	tok   Token  // last lexed token
	val   string // string value of last token (or "")
}

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

func parseSearchExpression(src []byte) (searchExpr *BinaryExpr, err error) {
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
	searchExpr = p.parse().(*BinaryExpr)

	return
}

func (p *parser) parse() Expr {
	var expr Expr

	for p.tok != EOL {
		switch p.tok {
		case AND, OR:
			if expr == nil {
				panic(p.errorf("expected expression at left of %s", p.tok))
			}
			op := p.tok
			p.next()
			p.expect(VAR_NAME)
			expr = &BinaryExpr{Left: expr, Op: op, Right: p.eqlExpr()}
		case VAR_NAME:
			expr = p.eqlExpr()
			if p.tok != EOL && !p.matches(AND, OR) {
				panic(p.errorf("expected operator after expression instead of %s", p.tok))
			}
		default:
			panic(p.errorf("expected expression instead of %s", p.tok))
		}
	}

	return expr
}

func (p *parser) eqlExpr() Expr {
	p.expect(VAR_NAME)
	name := p.val
	expr := &BinaryExpr{Left: &VarExpr{name}}

	p.next()

	switch p.tok {
	case GREATER, GTE, LESS, LTE, EQUALS, NOT_EQUALS:
		expr.Op = p.tok
	default:
		panic(p.errorf("expected operator instead of %s", p.tok))

	}

	p.next()
	switch p.tok {
	case STRING:
		expr.Right = &StrExpr{p.val}
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

// Format given string and args with Sprintf and return an error
// with that message and the current position.
func (p *parser) errorf(format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	return &ParseError{p.pos, message}
}
