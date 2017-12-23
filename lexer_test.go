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
		Token{Type: Exclamation},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"$", []Token{
		Token{Type: Dollar},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"(", []Token{
		Token{Type: OpenParen},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{")", []Token{
		Token{Type: ClosedParen},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"...", []Token{
		Token{Type: Spread, ColumnEnd: 2},
		Token{Type: EOF, ColumnStart: 3, ColumnEnd: 3}},
	},
	{":", []Token{
		Token{Type: Colon},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"=", []Token{
		Token{Type: Equals},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"@", []Token{
		Token{Type: At},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"[", []Token{
		Token{Type: OpenBracket},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"]", []Token{
		Token{Type: ClosedBracket},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"{", []Token{
		Token{Type: OpenBrace},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"}", []Token{
		Token{Type: ClosedBrace},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"|", []Token{
		Token{Type: VerticalBar},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},

	// Whitespace tests
	{" ", []Token{
		Token{Type: Whitespace},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"\u0009", []Token{
		Token{Type: Whitespace},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"\u0020", []Token{
		Token{Type: Whitespace},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"\t", []Token{
		Token{Type: Whitespace},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"\u0009 \u0020 \t", []Token{
		Token{Type: Whitespace, ColumnStart: 0, ColumnEnd: 0},
		Token{Type: Whitespace, ColumnStart: 1, ColumnEnd: 1},
		Token{Type: Whitespace, ColumnStart: 2, ColumnEnd: 2},
		Token{Type: Whitespace, ColumnStart: 3, ColumnEnd: 3},
		Token{Type: Whitespace, ColumnStart: 4, ColumnEnd: 4},
		Token{Type: EOF, ColumnStart: 5, ColumnEnd: 5}},
	},

	// Insignificant Comma (whitespace) tests
	{",", []Token{
		Token{Type: Whitespace},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{", , ,", []Token{
		Token{Type: Whitespace, ColumnStart: 0, ColumnEnd: 0},
		Token{Type: Whitespace, ColumnStart: 1, ColumnEnd: 1},
		Token{Type: Whitespace, ColumnStart: 2, ColumnEnd: 2},
		Token{Type: Whitespace, ColumnStart: 3, ColumnEnd: 3},
		Token{Type: Whitespace, ColumnStart: 4, ColumnEnd: 4},
		Token{Type: EOF, ColumnStart: 5, ColumnEnd: 5}},
	},

	// LineTerminator tests
	{"\u000A", []Token{
		Token{Type: LineTerminator},
		Token{Type: EOF, Line: 1}},
	},
	// Why borked?
	{"\u000D", []Token{
		Token{Type: LineTerminator},
		Token{Type: EOF, Line: 1}},
	},
	{"\u000D\u000A", []Token{
		Token{Type: LineTerminator, ColumnEnd: 1},
		Token{Type: EOF, Line: 1}},
	},
	{"\u000A\u000D\u000D\u000A", []Token{
		Token{Type: LineTerminator},
		Token{Type: LineTerminator, Line: 1},
		Token{Type: LineTerminator, ColumnEnd: 1, Line: 2},
		Token{Type: EOF, Line: 3}},
	},

	// Comment tests
	{"#", []Token{
		Token{Type: Comment},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"##", []Token{
		Token{Type: Comment, ColumnEnd: 1},
		Token{Type: EOF, ColumnStart: 2, ColumnEnd: 2}},
	},
	{"# This is a comment without a line terminator!", []Token{
		Token{Type: Comment, ColumnEnd: 45},
		Token{Type: EOF, ColumnStart: 46, ColumnEnd: 46}},
	},
	{"# This is a comment with a line terminator!\u000A", []Token{
		Token{Type: Comment, ColumnEnd: 42},
		Token{Type: LineTerminator, ColumnStart: 43, ColumnEnd: 43},
		Token{Type: EOF, Line: 1}},
	},
	{"# This is a comment with a line terminator!\u000D", []Token{
		Token{Type: Comment, ColumnEnd: 42},
		Token{Type: LineTerminator, ColumnStart: 43, ColumnEnd: 43},
		Token{Type: EOF, Line: 1}},
	},
	{"# This is a comment with a line terminator!\u000D\u000A", []Token{
		Token{Type: Comment, ColumnEnd: 42},
		Token{Type: LineTerminator, ColumnStart: 43, ColumnEnd: 44},
		Token{Type: EOF, Line: 1}},
	},
	{"##[](){} !$@=...:|,##", []Token{
		Token{Type: Comment, ColumnEnd: 20},
		Token{Type: EOF, ColumnStart: 21, ColumnEnd: 21}},
	},
	{"#~!@#$%^&*()_+1234567890-=qwertyuiop[]\\asdfghjkl;'zxvbnm,./QWERTYUIOP{}|ASDFGHJKL:\"ZXCVBNM<>?\t œ∑´®†¥¨ˆøπ“åß”åßf∆˚¬…˜æΩç√'\u000A", []Token{
		Token{Type: Comment, ColumnEnd: 121},
		Token{Type: LineTerminator, ColumnStart: 122, ColumnEnd: 122},
		Token{Type: EOF, Line: 1}},
	},

	// Name tests
	{"_", []Token{ // TODO Is just '_' a valid name? GraphQL spec seems to indicate it is, but seems wrong...
		Token{Type: Name, Value: "_"},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"a", []Token{
		Token{Type: Name, Value: "a"},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"A", []Token{
		Token{Type: Name, Value: "A"},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"_b", []Token{
		Token{Type: Name, Value: "_b", ColumnEnd: 1},
		Token{Type: EOF, ColumnStart: 2, ColumnEnd: 2}},
	},
	{"_B", []Token{
		Token{Type: Name, Value: "_B", ColumnEnd: 1},
		Token{Type: EOF, ColumnStart: 2, ColumnEnd: 2}},
	},
	{"_0", []Token{
		Token{Type: Name, Value: "_0", ColumnEnd: 1},
		Token{Type: EOF, ColumnStart: 2, ColumnEnd: 2}},
	},
	{"_1c", []Token{
		Token{Type: Name, Value: "_1c", ColumnEnd: 2},
		Token{Type: EOF, ColumnStart: 3, ColumnEnd: 3}},
	},
	{"abc", []Token{
		Token{Type: Name, Value: "abc", ColumnEnd: 2},
		Token{Type: EOF, ColumnStart: 3, ColumnEnd: 3}},
	},
	{"abc123", []Token{
		Token{Type: Name, Value: "abc123", ColumnEnd: 5},
		Token{Type: EOF, ColumnStart: 6, ColumnEnd: 6}},
	},
	{"_zyx987", []Token{
		Token{Type: Name, Value: "_zyx987", ColumnEnd: 6},
		Token{Type: EOF, ColumnStart: 7, ColumnEnd: 7}},
	},
	{"a_b_c_d", []Token{
		Token{Type: Name, Value: "a_b_c_d", ColumnEnd: 6},
		Token{Type: EOF, ColumnStart: 7, ColumnEnd: 7}},
	},
	{"_e_f_g_1_2_3_", []Token{
		Token{Type: Name, Value: "_e_f_g_1_2_3_", ColumnEnd: 12},
		Token{Type: EOF, ColumnStart: 13, ColumnEnd: 13}},
	},
	{"abcdefghjklmnopqrstuvwxyz0123456789", []Token{
		Token{Type: Name, Value: "abcdefghjklmnopqrstuvwxyz0123456789", ColumnEnd: 34},
		Token{Type: EOF, ColumnStart: 35, ColumnEnd: 35}},
	},
	{"ABCDEFGHJKLMNOPQRSTUVWXYZ0123456789", []Token{
		Token{Type: Name, Value: "ABCDEFGHJKLMNOPQRSTUVWXYZ0123456789", ColumnEnd: 34},
		Token{Type: EOF, ColumnStart: 35, ColumnEnd: 35}},
	},
	{"_abcdefghjklmnopqrstuvwxyz0123456789", []Token{
		Token{Type: Name, Value: "_abcdefghjklmnopqrstuvwxyz0123456789", ColumnEnd: 35},
		Token{Type: EOF, ColumnStart: 36, ColumnEnd: 36}},
	},
	{"_ABCDEFGHJKLMNOPQRSTUVWXYZ0123456789", []Token{
		Token{Type: Name, Value: "_ABCDEFGHJKLMNOPQRSTUVWXYZ0123456789", ColumnEnd: 35},
		Token{Type: EOF, ColumnStart: 36, ColumnEnd: 36}},
	},

	// String tests
	{"\"\"", []Token{
		Token{Type: String, Value: "", ColumnEnd: 1},
		Token{Type: EOF, ColumnStart: 2, ColumnEnd: 2}},
	},
	{"\"abc\"", []Token{
		Token{Type: String, Value: "abc", ColumnEnd: 4},
		Token{Type: EOF, ColumnStart: 5, ColumnEnd: 5}},
	},
	{"\"#[](){} !$@=...:|,\"", []Token{
		Token{Type: String, Value: "#[](){} !$@=...:|,", ColumnEnd: 19},
		Token{Type: EOF, ColumnStart: 20, ColumnEnd: 20}},
	},
	{"\"\u0020 \uFFFF\"", []Token{
		Token{Type: String, Value: "\u0020 \uFFFF", ColumnEnd: 4},
		Token{Type: EOF, ColumnStart: 5, ColumnEnd: 5}},
	},
	{"\"This is a long String with spaces, tabs\t punctuation, and smiles! :)\"", []Token{
		Token{Type: String, Value: "This is a long String with spaces, tabs\t punctuation, and smiles! :)", ColumnEnd: 69},
		Token{Type: EOF, ColumnStart: 70, ColumnEnd: 70}},
	},
	// TODO more string tests

	// Escaped character tests TODO more of them!
	{"\"\\b\"", []Token{
		Token{Type: String, Value: "\u0008", ColumnEnd: 3},
		Token{Type: EOF, ColumnStart: 4, ColumnEnd: 4}},
	},
	{"\"\\t\"", []Token{
		Token{Type: String, Value: "\u0009", ColumnEnd: 3},
		Token{Type: EOF, ColumnStart: 4, ColumnEnd: 4}},
	},
	{"\"\\n\"", []Token{
		Token{Type: String, Value: "\u000A", ColumnEnd: 3},
		Token{Type: EOF, ColumnStart: 4, ColumnEnd: 4}},
	},
	{"\"\\f\"", []Token{
		Token{Type: String, Value: "\u000C", ColumnEnd: 3},
		Token{Type: EOF, ColumnStart: 4, ColumnEnd: 4}},
	},
	{"\"\\r\"", []Token{
		Token{Type: String, Value: "\u000D", ColumnEnd: 3},
		Token{Type: EOF, ColumnStart: 4, ColumnEnd: 4}},
	},
	{"\"\\\"\"", []Token{
		Token{Type: String, Value: "\u0022", ColumnEnd: 3},
		Token{Type: EOF, ColumnStart: 4, ColumnEnd: 4}},
	},
	{"\"\\\\\"", []Token{
		Token{Type: String, Value: "\u005C", ColumnEnd: 3},
		Token{Type: EOF, ColumnStart: 4, ColumnEnd: 4}},
	},
	{"\"\\/\"", []Token{
		Token{Type: String, Value: "\u002F", ColumnEnd: 3},
		Token{Type: EOF, ColumnStart: 4, ColumnEnd: 4}},
	},

	// Escaped unicode tests
	{"\"\\u0000\"", []Token{
		Token{Type: String, Value: "\u0000", ColumnEnd: 7},
		Token{Type: EOF, ColumnStart: 8, ColumnEnd: 8}},
	},
	{"\"\\uFFFF\"", []Token{
		Token{Type: String, Value: "\uFFFF", ColumnEnd: 7},
		Token{Type: EOF, ColumnStart: 8, ColumnEnd: 8}},
	},
	{"\"\\u000A\"", []Token{
		Token{Type: String, Value: "\n", ColumnEnd: 7},
		Token{Type: EOF, ColumnStart: 8, ColumnEnd: 8}},
	},
	{"\"a\\u0062c\\u0064e\"", []Token{
		Token{Type: String, Value: "abcde", ColumnEnd: 16},
		Token{Type: EOF, ColumnStart: 17, ColumnEnd: 17}},
	},
	{"\"\\u0048\\u0065\\u006c\\u006c\\u006f\\u002c\\u0020\\u004b\\u0072\\u0069\\u0073\\u0074\\u0069\\u006e\\u0065\"", []Token{
		Token{Type: String, Value: "Hello, Kristine", ColumnEnd: 91},
		Token{Type: EOF, ColumnStart: 92, ColumnEnd: 92}},
	},

	// Integer tests
	{"0", []Token{
		Token{Type: Integer, Value: "0"},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"1", []Token{
		Token{Type: Integer, Value: "1"},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"5", []Token{
		Token{Type: Integer, Value: "5"},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"9", []Token{
		Token{Type: Integer, Value: "9"},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"-0", []Token{
		Token{Type: Integer, Value: "-0", ColumnEnd: 1},
		Token{Type: EOF, ColumnStart: 2, ColumnEnd: 2}},
	},
	{"-1", []Token{
		Token{Type: Integer, Value: "-1", ColumnEnd: 1},
		Token{Type: EOF, ColumnStart: 2, ColumnEnd: 2}},
	},
	{"-5", []Token{
		Token{Type: Integer, Value: "-5", ColumnEnd: 1},
		Token{Type: EOF, ColumnStart: 2, ColumnEnd: 2}},
	},
	{"-9", []Token{
		Token{Type: Integer, Value: "-9", ColumnEnd: 1},
		Token{Type: EOF, ColumnStart: 2, ColumnEnd: 2}},
	},
	{"1234567890", []Token{
		Token{Type: Integer, Value: "1234567890", ColumnEnd: 9},
		Token{Type: EOF, ColumnStart: 10, ColumnEnd: 10}},
	},
	{"-1234567890", []Token{
		Token{Type: Integer, Value: "-1234567890", ColumnEnd: 10},
		Token{Type: EOF, ColumnStart: 11, ColumnEnd: 11}},
	},
	{"123456789012345678901234567890123456789012345678901234567890", []Token{
		Token{Type: Integer, Value: "123456789012345678901234567890123456789012345678901234567890", ColumnEnd: 59},
		Token{Type: EOF, ColumnStart: 60, ColumnEnd: 60}},
	},

	// Float tests
	{"0.0", []Token{
		Token{Type: Float, Value: "0.0", ColumnEnd: 2},
		Token{Type: EOF, ColumnStart: 3, ColumnEnd: 3}},
	},
	{"1.0", []Token{
		Token{Type: Float, Value: "1.0", ColumnEnd: 2},
		Token{Type: EOF, ColumnStart: 3, ColumnEnd: 3}},
	},
	{"-1.0", []Token{
		Token{Type: Float, Value: "-1.0", ColumnEnd: 3},
		Token{Type: EOF, ColumnStart: 4, ColumnEnd: 4}},
	},
	{"-1.1", []Token{
		Token{Type: Float, Value: "-1.1", ColumnEnd: 3},
		Token{Type: EOF, ColumnStart: 4, ColumnEnd: 4}},
	},
	{"-1.0123456789", []Token{
		Token{Type: Float, Value: "-1.0123456789", ColumnEnd: 12},
		Token{Type: EOF, ColumnStart: 13, ColumnEnd: 13}},
	},
	{"1e0", []Token{
		Token{Type: Float, Value: "1e0", ColumnEnd: 2},
		Token{Type: EOF, ColumnStart: 3, ColumnEnd: 3}},
	},
	{"2e1", []Token{
		Token{Type: Float, Value: "2e1", ColumnEnd: 2},
		Token{Type: EOF, ColumnStart: 3, ColumnEnd: 3}},
	},
	{"1e23", []Token{
		Token{Type: Float, Value: "1e23", ColumnEnd: 3},
		Token{Type: EOF, ColumnStart: 4, ColumnEnd: 4}},
	},
	{"1E23", []Token{
		Token{Type: Float, Value: "1E23", ColumnEnd: 3},
		Token{Type: EOF, ColumnStart: 4, ColumnEnd: 4}},
	},
	{"123e45", []Token{
		Token{Type: Float, Value: "123e45", ColumnEnd: 5},
		Token{Type: EOF, ColumnStart: 6, ColumnEnd: 6}},
	},
	{"123E45", []Token{
		Token{Type: Float, Value: "123E45", ColumnEnd: 5},
		Token{Type: EOF, ColumnStart: 6, ColumnEnd: 6}},
	},
	{"1.1234567e89", []Token{
		Token{Type: Float, Value: "1.1234567e89", ColumnEnd: 11},
		Token{Type: EOF, ColumnStart: 12, ColumnEnd: 12}},
	},
	{"-1.1234567e89", []Token{
		Token{Type: Float, Value: "-1.1234567e89", ColumnEnd: 12},
		Token{Type: EOF, ColumnStart: 13, ColumnEnd: 13}},
	},
	{"1.1234567E89", []Token{
		Token{Type: Float, Value: "1.1234567E89", ColumnEnd: 11},
		Token{Type: EOF, ColumnStart: 12, ColumnEnd: 12}},
	},
	{"1.1234567e+89", []Token{
		Token{Type: Float, Value: "1.1234567e+89", ColumnEnd: 12},
		Token{Type: EOF, ColumnStart: 13, ColumnEnd: 13}},
	},
	{"1.1234567e-89", []Token{
		Token{Type: Float, Value: "1.1234567e-89", ColumnEnd: 12},
		Token{Type: EOF, ColumnStart: 13, ColumnEnd: 13}},
	},
	{"1.1234567E+89", []Token{
		Token{Type: Float, Value: "1.1234567E+89", ColumnEnd: 12},
		Token{Type: EOF, ColumnStart: 13, ColumnEnd: 13}},
	},
	{"1.1234567E-89", []Token{
		Token{Type: Float, Value: "1.1234567E-89", ColumnEnd: 12},
		Token{Type: EOF, ColumnStart: 13, ColumnEnd: 13}},
	},
	{"1e+23", []Token{
		Token{Type: Float, Value: "1e+23", ColumnEnd: 4},
		Token{Type: EOF, ColumnStart: 5, ColumnEnd: 5}},
	},
	{"1e-23", []Token{
		Token{Type: Float, Value: "1e-23", ColumnEnd: 4},
		Token{Type: EOF, ColumnStart: 5, ColumnEnd: 5}},
	},
	{"1E+23", []Token{
		Token{Type: Float, Value: "1E+23", ColumnEnd: 4},
		Token{Type: EOF, ColumnStart: 5, ColumnEnd: 5}},
	},
	{"1E-23", []Token{
		Token{Type: Float, Value: "1E-23", ColumnEnd: 4},
		Token{Type: EOF, ColumnStart: 5, ColumnEnd: 5}},
	},
	{"-1e+23", []Token{
		Token{Type: Float, Value: "-1e+23", ColumnEnd: 5},
		Token{Type: EOF, ColumnStart: 6, ColumnEnd: 6}},
	},
	// Possible BUG tests, following values are weird but allowed by the grammar
	{"1e01", []Token{
		Token{Type: Float, Value: "1e01", ColumnEnd: 3},
		Token{Type: EOF, ColumnStart: 4, ColumnEnd: 4}},
	},
	{"1e00", []Token{
		Token{Type: Float, Value: "1e00", ColumnEnd: 3},
		Token{Type: EOF, ColumnStart: 4, ColumnEnd: 4}},
	},
	{"1.234e00005", []Token{
		Token{Type: Float, Value: "1.234e00005", ColumnEnd: 10},
		Token{Type: EOF, ColumnStart: 11, ColumnEnd: 11}},
	},

	// TODO more string tests with escaped characters

	// UnicodeBOM tests
	{"\uFEFF", []Token{
		Token{Type: UnicodeBOM},
		Token{Type: EOF, ColumnStart: 1, ColumnEnd: 1}},
	},
	{"\uFEFF \uFEFF", []Token{
		Token{Type: UnicodeBOM},
		Token{Type: Whitespace, ColumnStart: 1, ColumnEnd: 1},
		Token{Type: UnicodeBOM, ColumnStart: 2, ColumnEnd: 2},
		Token{Type: EOF, ColumnStart: 3, ColumnEnd: 3}},
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

		actual, _ := Tokenize(strings.NewReader(test.input), false)

		if len(actual) < 1 {
			t.Errorf("Tokenize(%s): Returned no tokens", test.input)
		}

		if len(actual) != len(test.expected) {
			t.Errorf("Lexer(%s): expected %v, actual %v", test.input, test.expected, actual)
		} else {
			for i := 0; i < len(actual); i++ {
				if !reflect.DeepEqual(test.expected[i], actual[i]) {
					t.Errorf("Lexer(%s): expected %v, actual %v", test.input, test.expected, actual)
				}
			}
		}
	}

	// Negative test cases
	for _, test := range lexerTokenizeNegativeTests {
		_, actual := Tokenize(strings.NewReader(test.input), false)

		if actual == nil {
			t.Errorf("Tokenize(%s): Expected error, but was nil", test.input)
		}

		if !reflect.DeepEqual(test.expected, actual) {
			t.Errorf("Tokenize(%s): Expected error '%s', but got error '%s'", test.input, test.expected, actual)
		}
	}
}
