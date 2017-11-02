package graphql

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

// Tokenize tokenizes GraphQL documents from an io.Reader and returns
// a list of Tokens corressponding to the text. Returns an error if one occurs
// during the reading of runes from the Reader, or if there are invalid tokens
// in the text according to the GraphQL language specifcation
func Tokenize(r io.Reader) ([]Token, error) {
	var tokens []Token
	lexer := lexer{reader: bufio.NewReader(r)}

	for {
		token, err := lexer.nextToken()

		if err != nil {
			return tokens, err
		}
		tokens = append(tokens, token)

		if token.Type == EOF {
			return tokens, nil
		}
	}
}

// lexer represents a lexical token scanner for GraphQL
type lexer struct {
	reader *bufio.Reader
	// TODO Keep track of line & column indexes for errors, etc
}

// nextToken returns the next token and position from the underlying reader.
// Also returns the literal text read for strings, numbers, and duration tokens
// since these token types can have different literal representations.
func (l *lexer) nextToken() (Token, error) {
	var token Token

	r, _, err := l.reader.ReadRune()

	if err != nil {
		if err == io.EOF {
			return Token{Type: EOF, Value: ""}, nil
		}
		return Token{Type: Invalid}, err
	}

	switch {

	// Punctuators
	case r == '{':
		token = Token{Type: OpenBrace}
	case r == '}':
		token = Token{Type: ClosedBrace}
	case r == '(':
		token = Token{Type: OpenParen}
	case r == ')':
		token = Token{Type: ClosedParen}
	case r == '[':
		token = Token{Type: OpenBracket}
	case r == ']':
		token = Token{Type: ClosedBracket}
	case r == '=':
		token = Token{Type: Equals}
	case r == '!':
		token = Token{Type: Exclamation}
	case r == '$':
		token = Token{Type: Dollar}
	case r == ':':
		token = Token{Type: Colon}
	case r == '@':
		token = Token{Type: At}
	case r == '|':
		token = Token{Type: VerticalBar}

	// Spread "..."
	case r == '.':
		if s, err := l.consume(2, isPeriod); err == nil {
			token = Token{Type: Spread}
		} else {
			return Token{Type: Invalid}, fmt.Errorf("expected ... but found %c%s", r, s)
		}

	// Line Terminators (new line, carriage return, carriage return + new line)
	case r == '\u000D':
		// TODO test this error case thoroughly, if an error occurs during peek, a
		// LineTerminator can be returned, but the next call to nextToken must fail
		if nextr, _, err := l.peek(); err == nil {
			if nextr == '\u000A' {
				l.reader.ReadRune()
			}
		}
		fallthrough
	case r == '\u000A':
		token = Token{Type: LineTerminator}

	// Insignificant comma & whitespace (space, tab)
	case r == ',':
		fallthrough
	case r == '\u0009':
		fallthrough
	case r == '\u0020':
		token = Token{Type: Whitespace}

	// Name
	case r == '_':
		fallthrough
	case r >= 'a' && r <= 'z':
		fallthrough
	case r >= 'A' && r <= 'Z':
		if s, err := l.consumeAll(isNameCharacter); err == nil {
			token = Token{Type: Name, Value: fmt.Sprintf("%c%s", r, s)} // TODO SPrintf seems weird, rethink this?
		} else {
			return Token{Type: Invalid}, err
		}

	// String
	case r == '"':
		if t, err := l.consumeString(); err == nil {
			token = t
		} else {
			return Token{Type: Invalid}, err
		}

	// Integer or Float
	case r == '-':
		fallthrough
	case r >= '0' && r <= '9':
		if t, err := l.consumeNumber(r); err == nil {
			token = t
		} else {
			return Token{Type: Invalid}, err
		}

	// Comment
	case r == '#':
		if _, err := l.consumeAll(IsCommentCharacter); err == nil {
			token = Token{Type: Comment}
		} else {
			return Token{Type: Invalid}, err
		}

	// UnicodeBOM
	// TODO what if a UnicodeBOM appears after the first rune in the input
	case r == '\uFEFF':
		token = Token{Type: UnicodeBOM}

	// Invalid character for the start of a Token
	default:
		return Token{Type: Invalid}, fmt.Errorf("invalid character: %c", r)
	}

	return token, nil
}

// peek returns the next rune from the Reader without consuming it
func (l *lexer) peek() (rune, int, error) {
	r, size, err := l.reader.ReadRune()

	if err != nil {
		return r, size, err
	}
	l.reader.UnreadRune()
	return r, size, nil
}

// consume will attempt to read n number of contiguous runes that when passed to
// matches returns true. All read runes will be returned as a string. If there
// are not enough runes to read an error will be returned along with a string
// containing every rune read
func (l *lexer) consume(n int, matches func(rune) bool) (string, error) {
	var consumed bytes.Buffer

	for i := 0; i < n; i++ {
		r, _, err := l.reader.ReadRune()

		if err == nil {
			consumed.WriteRune(r)
			if !matches(r) {
				return consumed.String(), fmt.Errorf("consumed rune that does not match")
			}
		} else {
			return consumed.String(), err
		}
	}
	return consumed.String(), nil
}

// consumeAll will attempt to read contiguous runes that when passed to matches
// returns true. All read runes will be returned as a string. The last
// successfully read rune that does not match will be marked as Unread
func (l *lexer) consumeAll(matches func(rune) bool) (string, error) {
	var consumed bytes.Buffer

	for {
		r, _, err := l.reader.ReadRune()

		if err == nil {
			// Rune isn't a match, consider it unread and return everything else
			if !matches(r) {
				l.reader.UnreadRune()
				return consumed.String(), nil
			}
		} else {
			if err == io.EOF {
				return consumed.String(), nil
			}
			return consumed.String(), err
		}
		consumed.WriteRune(r)
	}
}

// consumeString attempts to consume a full string from lexer's Reader. Assumes
// the starting double quote has already been read
func (l *lexer) consumeString() (Token, error) {
	var value bytes.Buffer

	// consume all string characters until a closing quotation mark is found.
	// Since escape sequences may appear anywhere in the string they must be
	// parsed before the rest of the string is consumed and are considered to be
	// special string characters
	for {
		if s, err := l.consumeAll(isStringCharacter); err == nil {
			value.WriteString(s)

			if r, _, err := l.reader.ReadRune(); err == nil {
				switch r {
				// Final closing quotation; entire string has been consumed
				case '"':
					return Token{Type: String, Value: value.String()}, nil
				// Escape sequence
				case '\\':
					if err := l.consumeStringEscapeSequence(&value); err != nil {
						return Token{Type: Invalid}, err
					}
				// Invalid string input
				default:
					return Token{Type: Invalid}, fmt.Errorf("invalid String: %s%c", value.String(), r)
				}
			} else {
				return Token{Type: Invalid}, fmt.Errorf("invalid String: %s", value.String())
			}
		} else { // TODO don't just ignore errors from consumeAll/consume functions since they are likely reader errors
			return Token{Type: Invalid}, fmt.Errorf("invalid String: %s", value.String())
		}
	}
}

// consumeStringEscapeSequence consumes an escape character sequence or an
// escaped unicode character sequence from the lexer's Reader and writes the
// corressponding rune to the value buffer. Assumes the first forward Solidus
// has already been read from the lexer's reader
func (l *lexer) consumeStringEscapeSequence(value *bytes.Buffer) error {
	if r, _, err := l.reader.ReadRune(); err == nil {
		// Escaped Unicode character sequence
		if r == 'u' {
			// TODO does the unicode byte order matter here, or is it assumed
			unicodeHexString, err := l.consume(4, isHexidecimalCharacter)
			if err == nil {
				// Ignoring possible error from ParseInt because unicodeHexString has
				// already been checked for valid hex characters TODO CR & sanity check
				unicodeHexValue, _ := strconv.ParseInt(unicodeHexString, 16, 64)
				value.WriteRune(rune(unicodeHexValue))
				return nil
			}
			return fmt.Errorf("invalid escaped unicode value in String: %s%c%c%s", value.String(), '\\', r, unicodeHexString)
		}

		// Escaped character sequence
		if escapedCharacter, err := escapedCharacterToRune(r); err == nil {
			value.WriteRune(escapedCharacter)
			return nil
		}
		return fmt.Errorf("invalid escape sequence character in String: %s%c%c", value.String(), '\\', r)

	}
	return fmt.Errorf("invalid escape sequence in String: %s%c", value.String(), '\\')
}

// consumeNumber attempts to consume an integer or a float from the lexer's
// Reader. Assumes the first rune has already been read and is contextually
// important (unlike a String's first rune) so it is passed as an argument
func (l *lexer) consumeNumber(first rune) (Token, error) {
	var value bytes.Buffer
	value.WriteRune(first)

	// Type of Token to be returned; assume Integer until proven otherwise
	tokenType := Integer

	// If peek returns an error, or nextRune is not an integer character the
	// integer is valid if first does not indicate the beginning of a negative
	// number, otherwise consume it as a multi-character integer part
	if nextRune, _, err := l.peek(); err == nil && isIntegerCharacter(nextRune) {
		// Special case: if the integer part starts with -0 it can only be -0
		if first == '-' && nextRune == '0' {
			l.reader.ReadRune() // Discard the next rune which we know is 0
			value.WriteRune(nextRune)
		} else if s, err := l.consumeAll(isIntegerCharacter); err == nil {
			value.WriteString(s)
		}
	} else if first == '-' {
		return Token{Type: Invalid}, fmt.Errorf("invalid Integer: -")
	}

	// Fractional part
	if nextRune, _, err := l.peek(); err == nil {
		if nextRune == '.' {
			l.reader.ReadRune() // Discard the next rune which we know is .
			value.WriteRune(nextRune)

			if fractionalPart, err := l.consumeAll(isIntegerCharacter); err == nil && fractionalPart != "" {
				value.WriteString(fractionalPart)
			} else {
				return Token{Type: Invalid}, fmt.Errorf("invalid Float: %s%s", value.String(), fractionalPart)
			}

			// Number is considered a Float type since it contains a fractional part
			tokenType = Float
		}
	}

	// Exponent part
	if nextRune, _, err := l.peek(); err == nil {
		if nextRune == 'e' || nextRune == 'E' {
			l.reader.ReadRune() // Discard the next rune which we know is e or E
			value.WriteRune(nextRune)

			// Exponent indicator may be followed by + or -
			if nextRune, _, err := l.peek(); err == nil {
				if nextRune == '+' || nextRune == '-' {
					l.reader.ReadRune() // Discard the next rune which we know is + or -
					value.WriteRune(nextRune)
				}
			} else {
				// If peek failed the float is invalid since an exponent indicator must
				// be followed by at least one integer character
				return Token{Type: Invalid}, fmt.Errorf("invalid Float: %s", value.String())
			}

			// Exponent value
			// BUG this allows a number like 10e01 or 10e00, is this allowed?
			// The GraphQL spec describes an exponent indicator followed by a digit
			// list which includes 0 so it does seem to be allowed, but regular
			// integers cannot start with 0 and be followed by other values
			if exponentPart, err := l.consumeAll(isIntegerCharacter); err == nil && exponentPart != "" {
				value.WriteString(exponentPart)
			} else {
				return Token{Type: Invalid}, fmt.Errorf("invalid Float: %s%s", value.String(), exponentPart)
			}

			// Number is considered a Float type since it contains an exponent part
			tokenType = Float
		}
	}

	return Token{Type: tokenType, Value: value.String()}, nil
}

// isPeriod returns true if r == '.', false otherwise
func isPeriod(r rune) bool {
	return r == '.'
}

// IsCommentCharacter returns true if r is '\u0009' or between '\u0020-\uFFFF'
func IsCommentCharacter(r rune) bool {
	return r >= '\u0020' && r <= '\uFFFF' || r == '\u0009'
}

// isNameCharacter returns true if r is '_' or between 'a-z', 'A-Z', '0-9'
func isNameCharacter(r rune) bool {
	return r == '_' ||
		r >= 'a' && r <= 'z' ||
		r >= 'A' && r <= 'Z' ||
		r >= '0' && r <= '9'
}

// isStringCharacter returns true if r is
func isStringCharacter(r rune) bool {
	return r >= '\u0020' && r <= '\uFFFF' && r != '"' && r != '\\' ||
		r == '\u0009'
}

// escapedCharacterToRune returns the rune associated an escape sequence
func escapedCharacterToRune(c rune) (rune, error) {
	var r rune

	switch c {
	case 'b': // Backspace
		r = '\u0008'
	case 't': // Horizontal Tab
		r = '\u0009'
	case 'n': // Newline
		r = '\u000A'
	case 'f': // Formfeed
		r = '\u000C'
	case 'r': // Carriage Return
		r = '\u000D'
	case '"': // Double Quote
		r = '\u0022'
	case '\\': // Reverse Solidus (Back slash)
		r = '\u005C'
	case '/': // Forward Solidus (Forward slash)
		r = '\u002F'
	default:
		return rune(0), fmt.Errorf("invalid escape character %c", c)
	}

	return r, nil
}

// isHexidecimalCharacter returns true if r is between '0-9', 'a-f', 'A-F'
func isHexidecimalCharacter(r rune) bool {
	return r >= '0' && r <= '9' ||
		r >= 'a' && r <= 'f' ||
		r >= 'A' && r <= 'F'
}

// isIntegerCharacter returns true if r is between '0-9'. All integer are
// represented in Base 10
func isIntegerCharacter(r rune) bool {
	return r >= '0' && r <= '9'
}
