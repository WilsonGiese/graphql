package graphql

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

// Tokenize tokenizes GraphQL documents from an io.Reader and returns
// a list of Tokens corressponding to the text. Returns an error if one occurs
// during the reading of runes from the Reader, or if there are invalid tokens
// in the text according to the GraphQL language specifcation
func Tokenize(r io.Reader, ignoreWhitespace bool) ([]Token, error) {
	var tokens []Token
	lexer := lexer{reader: bufio.NewReader(r)}

	for {
		token, err := lexer.nextToken()

		if err != nil {
			return tokens, err
		}

		// Ignore Whitespace
		if ignoreWhitespace {
			if token.Type != Whitespace && token.Type != LineTerminator {
				tokens = append(tokens, token)
			}
		} else {
			tokens = append(tokens, token)
		}

		if token.Type == EOF {
			return tokens, nil
		}
	}
}

// lexer represents a lexical token scanner for GraphQL
type lexer struct {
	reader      *bufio.Reader
	line        int // Current reader line position
	column      int // Current reader column position
	savedColumn int // Previous reader column position before line increment
}

// Lexer errors
var errNotEnoughMatchingRunes = errors.New("not enough matching runes to read")

// nextToken returns the next token and position from the underlying reader.
// Also returns the literal text read for strings, numbers, and duration tokens
// since these token types can have different literal representations.
func (l *lexer) nextToken() (Token, error) {
	token := Token{Line: l.line, ColumnStart: l.column}

	r, _, err := l.readRune()

	if err != nil {
		if err == io.EOF {
			token.Type = EOF
			token.ColumnEnd = token.ColumnStart
			return token, nil
		}
		return InvalidToken, err
	}

	switch {

	// Punctuators
	case r == '{':
		token.Type = OpenBrace
		//token = tokeToken{Type: OpenBrace}
	case r == '}':
		token.Type = ClosedBrace
	case r == '(':
		token.Type = OpenParen
	case r == ')':
		token.Type = ClosedParen
	case r == '[':
		token.Type = OpenBracket
	case r == ']':
		token.Type = ClosedBracket
	case r == '=':
		token.Type = Equals
	case r == '!':
		token.Type = Exclamation
	case r == '$':
		token.Type = Dollar
	case r == ':':
		token.Type = Colon
	case r == '@':
		token.Type = At
	case r == '|':
		token.Type = VerticalBar

	// Spread "..."
	case r == '.':
		if s, err := l.consume(2, isPeriod); err == nil {
			token.Type = Spread
		} else {
			if err == io.EOF || err == errNotEnoughMatchingRunes {
				return InvalidToken, fmt.Errorf("expected ... but found %c%s", r, s)
			}
			return InvalidToken, err
		}

	// Line Terminators (new line, carriage return, carriage return + new line)
	case r == '\u000D':
		// TODO test this error case thoroughly, if an error occurs during peek, a
		// LineTerminator can be returned, but the next call to nextToken must fail
		if nextr, _, err := l.peek(); err == nil {
			if nextr == '\u000A' {
				l.readRune()
			}
		} else if err != io.EOF {
			return InvalidToken, err
		}
		fallthrough
	case r == '\u000A':
		token.Type = LineTerminator

	// Insignificant comma & whitespace (space, tab)
	case r == ',':
		fallthrough
	case r == '\u0009':
		fallthrough
	case r == '\u0020':
		token.Type = Whitespace

	// Name
	case r == '_':
		fallthrough
	case r >= 'a' && r <= 'z':
		fallthrough
	case r >= 'A' && r <= 'Z':
		if s, err := l.consumeAll(isNameCharacter); err == nil {
			token.Type = Name
			token.Value = fmt.Sprintf("%c%s", r, s)
		} else {
			return InvalidToken, err
		}

	// String
	case r == '"':
		if value, err := l.consumeString(); err == nil {
			token.Type = String
			token.Value = value
		} else {
			return InvalidToken, err
		}

	// Integer or Float
	case r == '-':
		fallthrough
	case r >= '0' && r <= '9':
		if value, t, err := l.consumeNumber(r); err == nil {
			token.Type = t
			token.Value = value
		} else {
			return InvalidToken, err
		}

	// Comment
	case r == '#':
		if _, err := l.consumeAll(IsCommentCharacter); err == nil {
			token.Type = Comment
		} else {
			return InvalidToken, err
		}

	// UnicodeBOM
	// TODO what if a UnicodeBOM appears after the first rune in the input
	case r == '\uFEFF':
		token.Type = UnicodeBOM
	// Invalid character for the start of a Token
	default:
		return InvalidToken, fmt.Errorf("invalid character: %c", r)
	}
	token.ColumnEnd = l.column - 1

	if token.Type == LineTerminator {
		l.incrementLine()
	}

	return token, nil
}

// peek returns the next rune from the Reader without consuming it
func (l *lexer) peek() (rune, int, error) {
	r, size, err := l.readRune()

	if err == nil {
		l.unreadRune()
		return r, size, nil
	}
	return r, size, err
}

// consume will attempt to read n number of contiguous runes that when passed to
// matches returns true. All read runes will be returned as a string. If there
// are not enough runes to read an error will be returned along with a string
// containing every rune read
func (l *lexer) consume(n int, matches func(rune) bool) (string, error) {
	var consumed bytes.Buffer

	for i := 0; i < n; i++ {
		if r, _, err := l.readRune(); err == nil {
			consumed.WriteRune(r)
			if !matches(r) {
				return consumed.String(), errNotEnoughMatchingRunes
			}
		} else {
			return consumed.String(), err
		}
	}
	return consumed.String(), nil
}

// consumeAll will attempt to read contiguous runes that when passed to matches
// returns true. All read runes will be returned as a string. The last
// successfully read rune that does not match will be marked as Unread.
func (l *lexer) consumeAll(matches func(rune) bool) (string, error) {
	var consumed bytes.Buffer

	for {
		if r, _, err := l.readRune(); err == nil {
			// If rune isn't a match, consider it unread and return everything else
			if !matches(r) {
				l.unreadRune()
				return consumed.String(), nil
			}
			consumed.WriteRune(r)
		} else {
			if err == io.EOF {
				return consumed.String(), nil
			}
			return consumed.String(), err
		}
	}
}

// consumeString attempts to consume a full string from lexer's Reader. Assumes
// the starting double quote has already been read
func (l *lexer) consumeString() (string, error) {
	var value bytes.Buffer

	// consume all string characters until a closing quotation mark is found.
	// Since escape sequences may appear anywhere in the string they must be
	// parsed before the rest of the string is consumed and are considered to be
	// special string characters
	for {
		if s, err := l.consumeAll(isStringCharacter); err == nil {
			value.WriteString(s)

			if r, _, err := l.readRune(); err == nil {
				switch r {
				// Final closing quotation; entire string has been consumed
				case '"':
					return value.String(), nil
				// Escape sequence
				case '\\':
					if err := l.consumeStringEscapeSequence(&value); err != nil {
						return value.String(), err
					}
				// Invalid string input
				default:
					return value.String(), fmt.Errorf("invalid String: %s%c", value.String(), r)
				}
			} else {
				if err == io.EOF {
					return value.String(), fmt.Errorf("invalid String: %s", value.String())
				}
				return value.String(), err
			}
		} else {
			return value.String(), err
		}
	}
}

// consumeStringEscapeSequence consumes an escape character sequence or an
// escaped unicode character sequence from the lexer's Reader and writes the
// corressponding rune to the value buffer. Assumes the first forward Solidus
// has already been read from the lexer's reader
func (l *lexer) consumeStringEscapeSequence(value *bytes.Buffer) error {
	r, _, err := l.readRune()

	if err == nil {
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

			if err == io.EOF || err == errNotEnoughMatchingRunes {
				return fmt.Errorf("invalid escaped unicode value in String: %s%c%c%s", value.String(), '\\', r, unicodeHexString)
			}
			return err
		}

		// Escaped character sequence
		if escapedCharacter, err := escapedCharacterToRune(r); err == nil {
			value.WriteRune(escapedCharacter)
			return nil
		}
		return fmt.Errorf("invalid escape sequence character in String: %s%c%c", value.String(), '\\', r)

	}

	if err == io.EOF {
		return fmt.Errorf("invalid escape sequence in String: %s%c", value.String(), '\\')
	}
	return err
}

// consumeNumber attempts to consume an integer or a float from the lexer's
// Reader. Assumes the first rune has already been read and is contextually
// important (unlike a String's first rune) so it is passed as an argument
func (l *lexer) consumeNumber(first rune) (string, TokenType, error) {
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
			l.readRune() // Discard the next rune which we know is 0
			value.WriteRune(nextRune)
		} else if s, err := l.consumeAll(isIntegerCharacter); err == nil {
			value.WriteString(s)
		} else {
			return value.String(), tokenType, err
		}
	} else if first == '-' {
		if err == io.EOF {
			if first == '-' {
				return value.String(), tokenType, fmt.Errorf("invalid Integer: -")
			}
			return value.String(), tokenType, nil
		}
		return value.String(), tokenType, err
	}

	// Fractional part
	if nextRune, _, err := l.peek(); err == nil {
		if nextRune == '.' {
			l.readRune() // Discard the next rune which we know is .
			value.WriteRune(nextRune)

			if fractionalPart, err := l.consumeAll(isIntegerCharacter); err == nil && fractionalPart != "" {
				value.WriteString(fractionalPart)
			} else {
				return value.String(), tokenType, fmt.Errorf("invalid Float: %s%s", value.String(), fractionalPart)
			}

			// Number is considered a Float type since it contains a fractional part
			tokenType = Float
		}
	} else if err != io.EOF {
		return value.String(), tokenType, err
	}

	// Exponent part
	if nextRune, _, err := l.peek(); err == nil {
		if nextRune == 'e' || nextRune == 'E' {
			l.readRune() // Discard the next rune which we know is e or E
			value.WriteRune(nextRune)

			// Exponent indicator may be followed by + or -
			if nextRune, _, err := l.peek(); err == nil {
				if nextRune == '+' || nextRune == '-' {
					l.readRune() // Discard the next rune which we know is + or -
					value.WriteRune(nextRune)
				}
			} else {
				if err == io.EOF {
					return value.String(), tokenType, fmt.Errorf("invalid Float: %s", value.String())
				}
				return value.String(), tokenType, err
			}

			// Exponent value
			// BUG this allows a number like 10e01 or 10e00, is this allowed?
			// The GraphQL spec describes an exponent indicator followed by a digit
			// list which includes 0 so it does seem to be allowed, but regular
			// integers cannot start with 0 and be followed by other values
			if exponentPart, err := l.consumeAll(isIntegerCharacter); err == nil {
				if exponentPart == "" {
					return value.String(), tokenType, fmt.Errorf("invalid Float: %s", value.String())
				}
				value.WriteString(exponentPart)
			} else {
				return value.String(), tokenType, err
			}

			// Number is considered a Float type since it contains an exponent part
			tokenType = Float
		}
	} else if err != io.EOF {
		return value.String(), tokenType, err
	}

	return value.String(), tokenType, nil
}

func (l *lexer) readRune() (rune, int, error) {
	l.incrementColumn()
	return l.reader.ReadRune()
}

func (l *lexer) unreadRune() error {
	l.undoLastIncrement()
	return l.reader.UnreadRune()
}

func (l *lexer) incrementLine() {
	l.line += 1
	l.savedColumn = l.column
	l.column = 0
}

func (l *lexer) incrementColumn() {
	l.column += 1
	l.savedColumn = -1
}

func (l *lexer) undoLastIncrement() {
	// If savedColumn is -1 it idicates the last increment was a column increment
	if l.savedColumn == -1 {
		l.column -= 1
	} else {
		l.line -= 1
		l.column = l.savedColumn
	}
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
