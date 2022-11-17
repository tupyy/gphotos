package filter

import (
	"strings"
	"testing"
)

func TestTokens(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{
			input:  "[ ] = != < <= > >= like name 'test' \"test\" ",
			output: "[ ] = != < <= > >= like variable string string EOL",
		},
		{
			input:  "name = 'test' in ",
			output: "variable = string in EOL",
		},
		{
			input:  "name = 'test' and description != 'toto' and location = 'loc' or",
			output: "variable = string and variable != string and variable = string or EOL",
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			l := newLexer([]byte(test.input))

			tokens := []string{}
			for {
				_, tok, _ := l.Scan()
				tokens = append(tokens, tok.String())
				if tok == EOL {
					break
				}

			}

			output := strings.Join(tokens, " ")
			if strings.TrimSpace(output) != test.output {
				t.Errorf("expected %q, got %q", test.output, output)
			}
		})
	}
}
