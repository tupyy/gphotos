package filter

import (
	"fmt"
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
			test:     "name = 'test' and description != 'toto' and location = 'loc'",
			expected: "(((\"name\" = \"test\") and (\"description\" != \"toto\")) and (\"location\" = \"loc\"))",
			hasError: false,
		},
		{
			test:     "name = 'test' and description != 'toto' and location = 'loc' or tag = 'tag'",
			expected: "((((\"name\" = \"test\") and (\"description\" != \"toto\")) and (\"location\" = \"loc\")) or (\"tag\" = \"tag\"))",
			hasError: false,
		},
		{
			test:     "name = 'test' or description != 'toto'",
			expected: "((\"name\" = \"test\") or (\"description\" != \"toto\"))",
			hasError: false,
		},
		{
			test:     "name = 'test' description != 'toto'",
			hasError: true,
		},
		{
			test:     "and name = 'test'",
			hasError: true,
		},
		{
			test:     "name = 'test' and",
			hasError: true,
		},
		{
			test:     "name and 'test'",
			hasError: true,
		},
		{
			test:     "name = 'test' =",
			hasError: true,
		},
		{
			test:     "name in ['1', '2']",
			expected: "(\"name\" in [1,2])",
			hasError: false,
		},
		{
			test:     "name in '1', '2']",
			hasError: true,
		},
		{
			test:     "name in ['1', '2'",
			hasError: true,
		},
		{
			test:     "name in ['1''2']",
			hasError: true,
		},
	}

	for idx, data := range exprs {
		t.Run(fmt.Sprintf("test%d: %s", idx+1, data.test), func(t *testing.T) {
			searchExpr, err := parse([]byte(data.test))
			if data.hasError {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, data.expected, searchExpr.String())
			}

		})
	}
}
