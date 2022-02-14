package search

import (
	"strings"
	"testing"
)

func TestKeywordToken(t *testing.T) {
	tests := []struct {
		name string
		tok  Token
	}{
		{"name", NAME},
		{"description", DESCRIPTION},
		{"location", LOCATION},
		{"date", DATE},
		{"owner", OWNER},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tok := keywordToken(test.name)
			if tok != test.tok {
				t.Errorf("expected %v, got %v", test.tok, tok)
			}
		})
	}
}

func TestAllTokens(t *testing.T) {
	input := "( ) = != ! < <= > >= name description location date owner tag 'test' \"test\" "

	strs := []string{}
	l := NewLexer([]byte(input))
	for {
		_, tok, _ := l.Scan()
		strs = append(strs, tok.String())
		if tok == EOL {
			break
		}
	}
	output := strings.Join(strs, " ")

	expected := "( ) = != ! < <= > >= name description location date owner tag string string EOL"

	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestValue(t *testing.T) {
	input := " 'test' \"test\" "

	strs := []string{}
	l := NewLexer([]byte(input))
	for {
		_, tok, val := l.Scan()
		strs = append(strs, val)
		if tok == EOL {
			break
		}
	}
	output := strings.Join(strs, " ")

	expected := "test test "

	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}

}
