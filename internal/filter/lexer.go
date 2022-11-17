package filter

import "strings"

type lexer struct {
	src     []byte
	ch      byte
	offset  int
	pos     int
	nextPos int
}

func newLexer(src []byte) *lexer {
	l := &lexer{src: src}
	l.next()

	return l
}

func (l *lexer) Scan() (int, Token, string) {
	for l.ch == ' ' || l.ch == '\t' {
		l.next()
	}

	if l.ch == 0 {
		return l.pos, EOL, ""
	}

	tok := ILLEGAL
	pos := l.pos
	val := ""

	ch := l.ch
	l.next()

	// keywords
	if isNameStart(ch) {
		start := l.offset - 2
		for isNameStart(l.ch) || isDot(l.ch) { // accept dots in the variable name like "user.permissions"
			l.next()
		}
		name := string(l.src[start : l.offset-1])
		switch strings.ToLower(name) {
		case "in":
			tok = IN
			val = ""
		case "and":
			tok = AND
		case "or":
			tok = OR
		case "like":
			tok = LIKE
		default:
			tok = VARIABLE
			val = name
		}
		return pos, tok, val
	}

	switch ch {
	case '[':
		tok = LBRACKET
	case ']':
		tok = RBRACKET
	case '=':
		tok = EQUALS
	case '!':
		switch l.ch {
		case '=':
			tok = NOT_EQUALS
			l.next()
		default:
			tok = ILLEGAL
		}
	case '<':
		switch l.ch {
		case '=':
			tok = LTE
			l.next()
		default:
			tok = LESS
		}
	case '>':
		switch l.ch {
		case '=':
			tok = GTE
			l.next()
		default:
			tok = GREATER
		}
	case ',':
		tok = COMMA
	case '"', '\'':
		chars := make([]byte, 0, 32) // most won't require heap allocation
		for l.ch != ch {
			c := l.ch
			if c == 0 {
				return l.pos, ILLEGAL, "didn't find end quote in string"
			}
			l.next()
			chars = append(chars, c)
		}
		l.next()
		tok = STRING
		val = string(chars)
	default:
		tok = ILLEGAL
		val = "unexpected char"
	}

	return l.pos, tok, val

}

// Load the next character into l.ch (or 0 on end of input) and update line position.
func (l *lexer) next() {
	l.pos = l.nextPos
	if l.offset >= len(l.src) {
		// For last character, move offset 1 past the end as it
		// simplifies offset calculations in NAME and NUMBER
		if l.ch != 0 {
			l.ch = 0
			l.offset++
			l.nextPos++
		}
		return
	}
	ch := l.src[l.offset]
	l.ch = ch
	l.nextPos++
	l.offset++
}

func isNameStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDot(ch byte) bool {
	return ch == '.'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
