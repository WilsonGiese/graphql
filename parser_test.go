package graphql

import "testing"

func TestParser(t *testing.T) {
	tokens := []Token{
		Token{Type: Name, Value: "query"},
		Token{Type: OpenBrace},
		Token{Type: ClosedBrace},
		Token{Type: EOF},
	}

	parser := Parser{tokens: tokens}

	err := parser.parse()
	if err != nil {
		t.Error(err)
	}
}
