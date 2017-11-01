package graphql

// Token represents a single GraphQL token from an set of characters
type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

// TokenType represents a lexical token for the GraphQL language
//go:generate stringer -type=TokenType
type TokenType int

// Lexical tokens for GraphQL
const (
	// Special tokens
	EOF TokenType = iota

	UnicodeBOM // Byte Order Mark (U+FEFF)

	// Horizontal Tab (U+0009)
	// Space (U+0020)
	Whitespace

	// New Line (U+000A)
	// Carriage Return, Lookahead != (U+000D)New Line (U+000A)
	// Carriage Return (U+000D), New Line (U+000A)
	LineTerminator

	// Symbols
	Comma // ,

	SourceCharacter // [\u0009\u000A\u000D\u0020-\uFFFF]

	Comment // # SoureCharacter until LineTerminator

	// Punctuators
	Exclamation   // !
	Dollar        // $
	OpenParen     // (
	ClosedParen   // )
	Spread        // ...
	Colon         // :
	Equals        // =
	At            // @
	OpenBracket   // [
	ClosedBracket // ]
	OpenBrace     // {
	ClosedBrace   // }
	VerticalBar   // |

	// [_A-Za-z][_0-9A-Za-z] ASCII only
	Name // e.g. 'abc'

	//
	Integer

	Float

	// ""
	// "[StringCharacters]" where StringCharacter is
	//    SourceCharacter but not " or \ or LineTerminator
	//    \u EscapedUnicode
	//    \  EscapedCharacter
	String

	// Invalid Token
	Invalid
)
