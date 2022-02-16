package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	exprs := []struct {
		test     string
		expected string
		hasError bool
	}{
		{
			test:     "name = 'test'",
			expected: "(\"name\" = \"test\")",
			hasError: false,
		},
		{
			test:     "name = 'test' & description != 'toto' & location = 'loc'",
			expected: "(((\"name\" = \"test\") & (\"description\" != \"toto\")) & (\"location\" = \"loc\"))",
			hasError: false,
		},
		{
			test:     "name = 'test' & description != 'toto' & location = 'loc' | tag = 'tag'",
			expected: "((((\"name\" = \"test\") & (\"description\" != \"toto\")) & (\"location\" = \"loc\")) | (\"tag\" = \"tag\"))",
			hasError: false,
		},
		{
			test:     "name = 'test' | description != 'toto'",
			expected: "((\"name\" = \"test\") | (\"description\" != \"toto\"))",
			hasError: false,
		},
		{
			test:     "name = 'test' description != 'toto'",
			hasError: true,
		},
		{
			test:     "& name = 'test'",
			hasError: true,
		},
		{
			test:     "name = 'test' &",
			hasError: true,
		},
		{
			test:     "name & 'test'",
			hasError: true,
		},
		{
			test:     "name = 'test' =",
			hasError: true,
		},
	}

	for _, data := range exprs {
		t.Run(data.test, func(t *testing.T) {
			searchExpr, err := parseSearchExpression([]byte(data.test))
			if data.hasError {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, searchExpr.String(), data.expected)
			}

		})
	}
}
