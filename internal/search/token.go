package search

type Token int

const (
	ILLEGAL Token = iota
	EOL

	AND
	EQUALS
	GTE
	GREATER
	LPAREN
	LTE
	LESS
	OR
	NOT_EQUALS
	RPAREN

	// literal names as (name, description, location..)
	VAR_NAME
	STRING
)

var tokenNames = map[Token]string{
	ILLEGAL:    "illegal",
	EOL:        "EOL",
	AND:        "&",
	VAR_NAME:   "name",
	EQUALS:     "=",
	GTE:        ">=",
	GREATER:    ">",
	LPAREN:     "(",
	LTE:        "<=",
	LESS:       "<",
	OR:         "|",
	RPAREN:     ")",
	NOT_EQUALS: "!=",
	STRING:     "string",
}

func (t Token) String() string {
	return tokenNames[t]
}
