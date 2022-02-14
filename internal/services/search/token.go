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
	NOT
	RPAREN

	// keywords
	NAME
	LOCATION
	DESCRIPTION
	OWNER
	DATE
	TAG

	STRING
)

var keywordTokens = map[string]Token{
	"name":        NAME,
	"date":        DATE,
	"description": DESCRIPTION,
	"owner":       OWNER,
	"tag":         TAG,
	"location":    LOCATION,
}

func keywordToken(keyword string) Token {
	return keywordTokens[keyword]
}

var tokenNames = map[Token]string{
	ILLEGAL:     "illegal",
	EOL:         "EOL",
	AND:         "&",
	NAME:        "name",
	DATE:        "date",
	DESCRIPTION: "description",
	EQUALS:      "=",
	GTE:         ">=",
	GREATER:     ">",
	LOCATION:    "location",
	LPAREN:      "(",
	LTE:         "<=",
	LESS:        "<",
	OR:          "|",
	OWNER:       "owner",
	RPAREN:      ")",
	TAG:         "tag",
	NOT_EQUALS:  "!=",
	NOT:         "!",
	STRING:      "string",
}

func (t Token) String() string {
	return tokenNames[t]
}
