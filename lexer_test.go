package graphql

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type PositiveLexerTokenizeTest struct {
	input    string
	expected []Token
}

type NegativeLexerTokenizeTest struct {
	input    string
	expected error
}

// Positive test case scenarios for Lexer.Tokenize
var lexerTokenizeTestsPositive = []PositiveLexerTokenizeTest{
	// Simple punctuator tests
	{"!", []Token{
		Token{Type: Exclamation}},
	},
	{"$", []Token{
		Token{Type: Dollar}},
	},
	{"(", []Token{
		Token{Type: OpenParen}},
	},
	{")", []Token{
		Token{Type: ClosedParen}},
	},
	{"...", []Token{
		Token{Type: Spread}},
	},
	{":", []Token{
		Token{Type: Colon}},
	},
	{"=", []Token{
		Token{Type: Equals}},
	},
	{"@", []Token{
		Token{Type: At}},
	},
	{"[", []Token{
		Token{Type: OpenBracket}},
	},
	{"]", []Token{
		Token{Type: ClosedBracket}},
	},
	{"{", []Token{
		Token{Type: OpenBrace}},
	},
	{"}", []Token{
		Token{Type: ClosedBrace}},
	},
	{"|", []Token{
		Token{Type: VerticalBar}},
	},

	// Whitespace tests
	{" ", []Token{
		Token{Type: Whitespace}},
	},
	{"\u0009", []Token{
		Token{Type: Whitespace}},
	},
	{"\u0020", []Token{
		Token{Type: Whitespace}},
	},
	{"\t", []Token{
		Token{Type: Whitespace}},
	},
	{"\u0009 \u0020 \t", []Token{
		Token{Type: Whitespace},
		Token{Type: Whitespace},
		Token{Type: Whitespace},
		Token{Type: Whitespace},
		Token{Type: Whitespace}},
	},

	// Insignificant Comma (whitespace) tests
	{",", []Token{
		Token{Type: Whitespace}},
	},
	{", , ,", []Token{
		Token{Type: Whitespace},
		Token{Type: Whitespace},
		Token{Type: Whitespace},
		Token{Type: Whitespace},
		Token{Type: Whitespace}},
	},

	// LineTerminator tests
	{"\u000A", []Token{
		Token{Type: LineTerminator}},
	},
	{"\u000D", []Token{
		Token{Type: LineTerminator}},
	},
	{"\u000D\u000A", []Token{
		Token{Type: LineTerminator}},
	},
	{"\u000A\u000D\u000D\u000A", []Token{
		Token{Type: LineTerminator},
		Token{Type: LineTerminator},
		Token{Type: LineTerminator}},
	},

	// Comment tests
	{"#", []Token{
		Token{Type: Comment}},
	},
	{"##", []Token{
		Token{Type: Comment}},
	},
	{"# This is a comment without a line terminator!", []Token{
		Token{Type: Comment}},
	},
	{"# This is a comment with a line terminator!\u000A", []Token{
		Token{Type: Comment},
		Token{Type: LineTerminator}},
	},
	{"# This is a comment with a line terminator!\u000D", []Token{
		Token{Type: Comment},
		Token{Type: LineTerminator}},
	},
	{"# This is a comment with a line terminator!\u000D\u000A", []Token{
		Token{Type: Comment},
		Token{Type: LineTerminator}},
	},
	{"##[](){} !$@=...:|,##", []Token{
		Token{Type: Comment}},
	},
	{"#~!@#$%^&*()_+1234567890-=qwertyuiop[]\\asdfghjkl;'zxvbnm,./QWERTYUIOP{}|ASDFGHJKL:\"ZXCVBNM<>?\t œ∑´®†¥¨ˆøπ“åß”åßf∆˚¬…˜æΩç√'\u000A", []Token{
		Token{Type: Comment},
		Token{Type: LineTerminator}},
	},

	// Name tests
	{"_", []Token{ // TODO Is just '_' a valid name? GraphQL spec seems to indicate it is, but seems wrong...
		Token{Type: Name, Value: "_"}},
	},
	{"a", []Token{
		Token{Type: Name, Value: "a"}},
	},
	{"A", []Token{
		Token{Type: Name, Value: "A"}},
	},
	{"_b", []Token{
		Token{Type: Name, Value: "_b"}},
	},
	{"_B", []Token{
		Token{Type: Name, Value: "_B"}},
	},
	{"_0", []Token{
		Token{Type: Name, Value: "_0"}},
	},
	{"_1c", []Token{
		Token{Type: Name, Value: "_1c"}},
	},
	{"abc", []Token{
		Token{Type: Name, Value: "abc"}},
	},
	{"abc123", []Token{
		Token{Type: Name, Value: "abc123"}},
	},
	{"_zyx987", []Token{
		Token{Type: Name, Value: "_zyx987"}},
	},
	{"a_b_c_d", []Token{
		Token{Type: Name, Value: "a_b_c_d"}},
	},
	{"_e_f_g_1_2_3_", []Token{
		Token{Type: Name, Value: "_e_f_g_1_2_3_"}},
	},
	{"abcdefghjklmnopqrstuvwxyz0123456789", []Token{
		Token{Type: Name, Value: "abcdefghjklmnopqrstuvwxyz0123456789"}},
	},
	{"ABCDEFGHJKLMNOPQRSTUVWXYZ0123456789", []Token{
		Token{Type: Name, Value: "ABCDEFGHJKLMNOPQRSTUVWXYZ0123456789"}},
	},
	{"_abcdefghjklmnopqrstuvwxyz0123456789", []Token{
		Token{Type: Name, Value: "_abcdefghjklmnopqrstuvwxyz0123456789"}},
	},
	{"_ABCDEFGHJKLMNOPQRSTUVWXYZ0123456789", []Token{
		Token{Type: Name, Value: "_ABCDEFGHJKLMNOPQRSTUVWXYZ0123456789"}},
	},

	// String tests
	{"\"\"", []Token{
		Token{Type: String, Value: ""}},
	},
	{"\"abc\"", []Token{
		Token{Type: String, Value: "abc"}},
	},
	{"\"#[](){} !$@=...:|,\"", []Token{
		Token{Type: String, Value: "#[](){} !$@=...:|,"}},
	},
	{"\"\u0020 \uFFFF\"", []Token{
		Token{Type: String, Value: "\u0020 \uFFFF"}},
	},
	{"\"This is a long String with spaces, tabs\t punctuation, and smiles! :)\"", []Token{
		Token{Type: String, Value: "This is a long String with spaces, tabs\t punctuation, and smiles! :)"}},
	},
	// TODO more string tests

	// Escaped character tests TODO more of them!
	{"\"\\b\"", []Token{
		Token{Type: String, Value: "\u0008"}},
	},
	{"\"\\t\"", []Token{
		Token{Type: String, Value: "\u0009"}},
	},
	{"\"\\n\"", []Token{
		Token{Type: String, Value: "\u000A"}},
	},
	{"\"\\f\"", []Token{
		Token{Type: String, Value: "\u000C"}},
	},
	{"\"\\r\"", []Token{
		Token{Type: String, Value: "\u000D"}},
	},
	{"\"\\\"\"", []Token{
		Token{Type: String, Value: "\u0022"}},
	},
	{"\"\\\\\"", []Token{
		Token{Type: String, Value: "\u005C"}},
	},
	{"\"\\/\"", []Token{
		Token{Type: String, Value: "\u002F"}},
	},

	// Escaped unicode tests
	{"\"\\u0000\"", []Token{
		Token{Type: String, Value: "\u0000"}},
	},
	{"\"\\uFFFF\"", []Token{
		Token{Type: String, Value: "\uFFFF"}},
	},
	{"\"\\u000A\"", []Token{
		Token{Type: String, Value: "\n"}},
	},
	{"\"a\\u0062c\\u0064e\"", []Token{
		Token{Type: String, Value: "abcde"}},
	},
	{"\"\\u0048\\u0065\\u006c\\u006c\\u006f\\u002c\\u0020\\u004b\\u0072\\u0069\\u0073\\u0074\\u0069\\u006e\\u0065\"", []Token{
		Token{Type: String, Value: "Hello, Kristine"}},
	},

	// Integer tests
	{"0", []Token{
		Token{Type: Integer, Value: "0"}},
	},
	{"1", []Token{
		Token{Type: Integer, Value: "1"}},
	},
	{"5", []Token{
		Token{Type: Integer, Value: "5"}},
	},
	{"9", []Token{
		Token{Type: Integer, Value: "9"}},
	},
	{"-0", []Token{
		Token{Type: Integer, Value: "-0"}},
	},
	{"-1", []Token{
		Token{Type: Integer, Value: "-1"}},
	},
	{"-5", []Token{
		Token{Type: Integer, Value: "-5"}},
	},
	{"-9", []Token{
		Token{Type: Integer, Value: "-9"}},
	},
	{"1234567890", []Token{
		Token{Type: Integer, Value: "1234567890"}},
	},
	{"-1234567890", []Token{
		Token{Type: Integer, Value: "-1234567890"}},
	},
	{"123456789012345678901234567890123456789012345678901234567890", []Token{
		Token{Type: Integer, Value: "123456789012345678901234567890123456789012345678901234567890"}},
	},

	// Float tests
	{"0.0", []Token{
		Token{Type: Float, Value: "0.0"}},
	},
	{"1.0", []Token{
		Token{Type: Float, Value: "1.0"}},
	},
	{"-1.0", []Token{
		Token{Type: Float, Value: "-1.0"}},
	},
	{"-1.1", []Token{
		Token{Type: Float, Value: "-1.1"}},
	},
	{"-1.0123456789", []Token{
		Token{Type: Float, Value: "-1.0123456789"}},
	},
	{"1e0", []Token{
		Token{Type: Float, Value: "1e0"}},
	},
	{"2e1", []Token{
		Token{Type: Float, Value: "2e1"}},
	},
	{"1e23", []Token{
		Token{Type: Float, Value: "1e23"}},
	},
	{"1E23", []Token{
		Token{Type: Float, Value: "1E23"}},
	},
	{"123e45", []Token{
		Token{Type: Float, Value: "123e45"}},
	},
	{"123E45", []Token{
		Token{Type: Float, Value: "123E45"}},
	},
	{"1.1234567e89", []Token{
		Token{Type: Float, Value: "1.1234567e89"}},
	},
	{"-1.1234567e89", []Token{
		Token{Type: Float, Value: "-1.1234567e89"}},
	},
	{"1.1234567E89", []Token{
		Token{Type: Float, Value: "1.1234567E89"}},
	},
	{"1.1234567e+89", []Token{
		Token{Type: Float, Value: "1.1234567e+89"}},
	},
	{"1.1234567e-89", []Token{
		Token{Type: Float, Value: "1.1234567e-89"}},
	},
	{"1.1234567E+89", []Token{
		Token{Type: Float, Value: "1.1234567E+89"}},
	},
	{"1.1234567E-89", []Token{
		Token{Type: Float, Value: "1.1234567E-89"}},
	},
	{"1e+23", []Token{
		Token{Type: Float, Value: "1e+23"}},
	},
	{"1e-23", []Token{
		Token{Type: Float, Value: "1e-23"}},
	},
	{"1E+23", []Token{
		Token{Type: Float, Value: "1E+23"}},
	},
	{"1E-23", []Token{
		Token{Type: Float, Value: "1E-23"}},
	},
	{"-1e+23", []Token{
		Token{Type: Float, Value: "-1e+23"}},
	},
	// Possible BUG tests, following values are weird but allowed by the grammar
	{"1e01", []Token{
		Token{Type: Float, Value: "1e01"}},
	},
	{"1e00", []Token{
		Token{Type: Float, Value: "1e00"}},
	},
	{"1.234e00005", []Token{
		Token{Type: Float, Value: "1.234e00005"}},
	},

	// TODO more string tests with escaped characters

	// UnicodeBOM tests
	{"\uFEFF", []Token{
		Token{Type: UnicodeBOM}},
	},
	{"\uFEFF \uFEFF", []Token{
		Token{Type: UnicodeBOM},
		Token{Type: Whitespace},
		Token{Type: UnicodeBOM}},
	},
}

var lexerTokenizeNegativeTests = []NegativeLexerTokenizeTest{
	// Spread tests
	{".", fmt.Errorf("expected ... but found .")},
	{"..", fmt.Errorf("expected ... but found ..")},
	{"....", fmt.Errorf("expected ... but found .")},
	{"..>", fmt.Errorf("expected ... but found ..>")},

	// String tests
	{"\"", fmt.Errorf("invalid String: ")},
	{"\"a", fmt.Errorf("invalid String: a")},
	{"\"Hello, World", fmt.Errorf("invalid String: Hello, World")},
	{"\"This is a string without an ending \\\"", fmt.Errorf("invalid String: This is a string without an ending \"")},
	{"\"This is a string with an invalid string character \u0000\"", fmt.Errorf("invalid String: This is a string with an invalid string character \u0000")},

	// Escape Sequence tests
	{"\"\\", fmt.Errorf("invalid escape sequence in String: \\")},
	{"\"\\ \"", fmt.Errorf("invalid escape sequence character in String: \\ ")},
	{"\"\\a\"", fmt.Errorf("invalid escape sequence character in String: \\a")},
	{"\"\\a\\b\"", fmt.Errorf("invalid escape sequence character in String: \\a")},
	{"\"some stuff\\cother stuff\"", fmt.Errorf("invalid escape sequence character in String: some stuff\\c")},

	{"\"\\u\"", fmt.Errorf("invalid escaped unicode value in String: \\u\"")},
	{"\"\\u0\"", fmt.Errorf("invalid escaped unicode value in String: \\u0\"")},
	{"\"\\u12\"", fmt.Errorf("invalid escaped unicode value in String: \\u12\"")},
	{"\"\\u34A\"", fmt.Errorf("invalid escaped unicode value in String: \\u34A\"")},
	{"\"\\uG\"", fmt.Errorf("invalid escaped unicode value in String: \\uG")},
	{"\"\\u123x\"", fmt.Errorf("invalid escaped unicode value in String: \\u123x")},

	// Invalid token tests (starting characters)
	{"?", fmt.Errorf("invalid character: ?")},
	{"^", fmt.Errorf("invalid character: ^")},
	{"\u0000", fmt.Errorf("invalid character: \u0000")},
	{"\uFFFF", fmt.Errorf("invalid character: \uFFFF")},
	{"abc 123 \uFFFF", fmt.Errorf("invalid character: \uFFFF")},
	{"ilike\uBEEFdoyou?", fmt.Errorf("invalid character: \uBEEF")},
	{"noiamv\u000Eg\u000An", fmt.Errorf("invalid character: \u000E")},

	// Integer tests
	{"-", fmt.Errorf("invalid Integer: -")},

	// Float tests
	{"1.", fmt.Errorf("invalid Float: 1.")},
	{"-1.", fmt.Errorf("invalid Float: -1.")},
	{"5.", fmt.Errorf("invalid Float: 5.")},
	{"-5.", fmt.Errorf("invalid Float: -5.")},
	{"2e", fmt.Errorf("invalid Float: 2e")},
	{"-2e", fmt.Errorf("invalid Float: -2e")},
	{"3.1e", fmt.Errorf("invalid Float: 3.1e")},
	{"-3.1e", fmt.Errorf("invalid Float: -3.1e")},
	{"4.2E", fmt.Errorf("invalid Float: 4.2E")},
	{"5.3e+", fmt.Errorf("invalid Float: 5.3e+")},
	{"56.34E+", fmt.Errorf("invalid Float: 56.34E+")},
	{"78.45e-", fmt.Errorf("invalid Float: 78.45e-")},
	{"87.56E-", fmt.Errorf("invalid Float: 87.56E-")},
}

// TestLexer preforms tests in the Lexer.Tokenize test tables:
// lexerTokenizePositiveTests
// lexerTokenizeNegativeTests
func TestLexerTokenize(t *testing.T) {

	// Positive test cases
	for _, test := range lexerTokenizeTestsPositive {

		actual, _ := Tokenize(strings.NewReader(test.input))

		if len(actual) < 1 {
			t.Errorf("Tokenize(%s): Returned no tokens", test.input)
		}

		for i := 0; i < len(actual)-1; i++ {
			if !reflect.DeepEqual(test.expected[i], actual[i]) {
				t.Errorf("Lexer(%s): expected %v, actual %v", test.input, test.expected, actual)
			}
		}

		if actual[len(actual)-1].Type != EOF {
			t.Errorf("Tokenize(%s): Final token type was not EOF, but was %s", test.input, actual[len(actual)-1].Type)
		}
	}

	// Negative test cases
	for _, test := range lexerTokenizeNegativeTests {
		_, actual := Tokenize(strings.NewReader(test.input))

		if actual == nil {
			t.Errorf("Tokenize(%s): Expected error, but was nil", test.input)
		}

		if !reflect.DeepEqual(test.expected, actual) {
			t.Errorf("Tokenize(%s): Expected error '%s', but got error '%s'", test.input, test.expected, actual)
		}
	}
}
