package search

import (
	"strings"
	"testing"
)

func TestAllTokens(t *testing.T) {
	input := "( ) = != < <= > >= name 'test' \"test\" "

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

	expected := "( ) = != < <= > >= name string string EOL"

	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestExpr(t *testing.T) {
	input := "name = 'test'"

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

	expected := "name = string EOL"

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
