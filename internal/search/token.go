package search

type Token int

const (
	ILLEGAL Token = iota
	EOL

	AND
	EQUALS
	DIV
	GTE
	GREATER
	LPAREN
	LTE
	LESS
	OR
	TILDA
	NOT_EQUALS
	RPAREN

	// literal names as (name, description, location..)
	VAR_NAME
	STRING
	DATE
	REGEX // regex are defined as /regex/
)

var tokenNames = map[Token]string{
	ILLEGAL:    "illegal",
	EOL:        "EOL",
	AND:        "&",
	DIV:        "/",
	EQUALS:     "=",
	GTE:        ">=",
	GREATER:    ">",
	LPAREN:     "(",
	LTE:        "<=",
	LESS:       "<",
	OR:         "|",
	TILDA:      "~",
	RPAREN:     ")",
	NOT_EQUALS: "!=",
	VAR_NAME:   "name",
	STRING:     "string",
	DATE:       "date",
	REGEX:      "regex",
}

func (t Token) String() string {
	return tokenNames[t]
}
