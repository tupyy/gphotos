package search

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokens(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{
			input:  "( ) = != < <= > >= ~ name 'test' \"test\" ",
			output: "( ) = != < <= > >= ~ name string string EOL",
		},
		{
			input:  "name = 'test'",
			output: "name = string EOL",
		},
		{
			input:  "name = 'test' & description != 'toto' & location = 'loc'",
			output: "name = string & name != string & name = string EOL",
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			l := NewLexer([]byte(test.input))

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

func TestValues(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{
			input:  "'test' \"test\" 02/02/2002 20/11/2002 /abc/",
			output: "test test 02/02/2002 20/11/2002 abc",
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			l := NewLexer([]byte(test.input))

			tokens := []string{}
			for {
				_, tok, val := l.Scan()
				tokens = append(tokens, val)
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

func TestErrors(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{
			input:  "'test test",
			output: "didn't find end quote in string",
		},
		{
			input:  "/ss",
			output: "didn't find / at the end of regex",
		},
		{
			input:  "/ss",
			output: "didn't find / at the end of regex",
		},
		{
			input:  "20/20/22",
			output: "expected 4 digits for year",
		},
		{
			input:  "20/20/22222",
			output: "expected 4 digits for year",
		},
		{
			input:  "20/220/2222",
			output: "expected 2 digits for month",
		},
		{
			input:  "20/2/2222",
			output: "expected 2 digits for month",
		},
		{
			input:  "2/22/2222",
			output: "expected 2 digits for day",
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			l := NewLexer([]byte(test.input))

			values := []string{}
			for {
				_, tok, val := l.Scan()
				values = append(values, val)
				if tok == EOL || tok == ILLEGAL {
					break
				}
			}

			output := strings.Join(values, " ")
			assert.Equal(t, test.output, strings.TrimSpace(output))
		})
	}
}
