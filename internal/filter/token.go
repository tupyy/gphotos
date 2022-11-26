package filter

type Token int

const (
	ILLEGAL Token = iota
	EOL

	AND
	EQUALS
	GTE
	GREATER
	LTE
	LESS
	OR
	NOT_EQUALS
	LIKE
	IN
	LBRACKET
	RBRACKET
	COMMA

	// literal names as (name, description, location..)
	STRING
	VARIABLE
)

var tokenNames = map[Token]string{
	ILLEGAL:    "illegal",
	EOL:        "EOL",
	AND:        "and",
	EQUALS:     "=",
	GTE:        ">=",
	GREATER:    ">",
	LTE:        "<=",
	LESS:       "<",
	OR:         "or",
	NOT_EQUALS: "!=",
	STRING:     "string",
	LBRACKET:   "[",
	RBRACKET:   "]",
	COMMA:      ",",
	IN:         "in",
	VARIABLE:   "variable",
	LIKE:       "like",
}

func (t Token) String() string {
	return tokenNames[t]
}
