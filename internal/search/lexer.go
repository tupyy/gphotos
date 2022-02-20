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
		tok := VAR_NAME
		val = name

		return pos, tok, val
	}

	switch ch {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// this has to be a date in format 01/02/2020
		start := l.offset - 2
		// check the day
		countDigits := 1
		for isDigit(l.ch) {
			countDigits++
			l.next()
		}
		if l.ch != '/' || countDigits != 2 {
			return l.pos, ILLEGAL, "expected 2 digits for day"
		}

		l.next()
		countDigits = 0
		for isDigit(l.ch) {
			countDigits++
			l.next()
		}
		if l.ch != '/' || countDigits != 2 {
			return l.pos, ILLEGAL, "expected 2 digits for month"
		}

		l.next()
		countDigits = 0
		for isDigit(l.ch) {
			countDigits++
			l.next()
		}
		if countDigits != 4 {
			return l.pos, ILLEGAL, "expected 4 digits for year"
		}

		tok = DATE
		val = string(l.src[start : l.offset-1])
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
	case '/':
		chars := make([]byte, 0, 32)
		for l.ch != ch {
			c := l.ch
			if c == 0 {
				return l.pos, ILLEGAL, "didn't find / at the end of regex"
			}
			l.next()
			chars = append(chars, c)
		}
		l.next()
		tok = REGEX
		val = string(chars)
	case '~':
		switch l.ch {
		case '=':
			tok = TILDA_EQUALS
			l.next()
		default:
			tok = ILLEGAL
			val = "unexpected character after ~"
		}
	default:
		tok = ILLEGAL
		val = "unexpected char"
	}

	return l.pos, tok, val

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

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
