package search

type Lexer struct {
	src     []byte
	ch      byte
	offset  int
	pos     int
	nextPos int
}

func NewLexer(src []byte) *Lexer {
	l := &Lexer{src: src}
	l.next()

	return l
}

func (l *Lexer) Scan() (int, Token, string) {
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
		for isNameStart(l.ch) {
			l.next()
		}
		name := string(l.src[start : l.offset-1])
		tok := keywordToken(name)

		return pos, tok, val
	}

	switch ch {
	case '(':
		tok = LPAREN
	case ')':
		tok = RPAREN
	case '=':
		tok = EQUALS
	case '!':
		switch l.ch {
		case '=':
			tok = NOT_EQUALS
			l.next()
		default:
			tok = NOT
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
	case '&':
		tok = AND
	case '|':
		tok = OR
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

	return pos, tok, val

}

// Load the next character into l.ch (or 0 on end of input) and update line position.
func (l *Lexer) next() {
	l.pos = l.nextPos
	if l.offset >= len(l.src) {
		// For last character, move offset 1 past the end as it
		// simplifies offset calculations in NAME and NUMBER
		if l.ch != 0 {
			l.ch = 0
			l.offset++
		}
		return
	}
	ch := l.src[l.offset]
	l.ch = ch
	l.offset++
}

func isNameStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}
