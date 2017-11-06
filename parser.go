package graphql

import (
	"fmt"
	"strings"
)

// Parser for GraphQL
type Parser struct {
	tokens   []Token
	position int
}

func (p *Parser) parse() (err error) {

	// Defer function to recover from a panic during parsing, and finally return
	// the panic back as a proper error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	p.parseDocument()

	return
}

func (p *Parser) parseDocument() Document {
	// A document starts with an OpenBrace if it is using query short-hand syntax,
	// otherwise a document is a list of Operations with their type stated
	if p.peek().Type == OpenBrace {

		p.parseSelectionSet()
		p.expect(EOF)
	} else {
		var operations []Operation
		var fragments []Fragment

		for {
			token := p.peek()
			if token.Type == EOF {
				break
			}
			defintionType := p.accept(Name, "query", "mutation", "fragment").Value

			if defintionType == "fragment" {
				fragments = append(fragments, p.parseFragment())
			} else {
				operations = append(operations, p.parseOperation(defintionType))
			}
		}
	}

	return Document{}
}

// OperationDefinition
// OperationType Name(opt) VariableDefinitions(opt) Directives(opt) SelectionSet
// SelectionSet
func (p *Parser) parseOperation(operationType string) Operation {
	operationName, _ := p.optional(Name)
	variables := p.parseVariables()
	directives := p.parseDirectives()
	selectionSet := p.parseSelectionSet()

	return Operation{
		operationType: operationType,
		operationName: operationName.Value,
		variables:     variables,
		directives:    directives,
		selectionSet:  selectionSet,
	}
}

func (p *Parser) parseFragment() Fragment {
	return Fragment{}
}

func (p *Parser) parseVariables() (variables []Variable) {

	p.expect(OpenParen)
	for {
		token := p.peek()

		if token.Type != Dollar {
			break
		}

		p.expect(Dollar)
		variableName := p.expect(Name)
		p.expect(Colon)
		variableType := p.parseType()

		_, defaultGiven := p.optional(Equals)

		if defaultGiven {
			switch p.peek().Type {
			case Name:
			case String:
			case Integer:
			case Float:
			}
		}

		//variableDefaultValue := p.parseValue()

	}
	p.expect(ClosedParen)

	return
}

func (p *Parser) parseType() Type {
	t := p._parseType()

	_, nonNull := p.optional(Exclamation)
	t.NonNull = nonNull

	return t
}

func (p *Parser) _parseType() (t Type) {
	token := p.peek()

	if token.Type == OpenBracket {
		p.expect(OpenBracket)
		subType := p.parseType()
		p.expect(ClosedBracket)

		t.List = true
		t.SubType = &subType
	} else {
		t.Type = p.expect(Name).Value
	}

	return
}

func (p *Parser) parseDirectives() Directives {
	return Directives{}
}

func (p *Parser) parseSelectionSet() SelectionSet {
	p.expect(OpenBrace)
	p.expect(ClosedBrace)
	return SelectionSet{}
}

func (p *Parser) accept(t TokenType, values ...string) Token {
	token := p.peek()

	if token.Type == t {
		for _, value := range values {
			if value == token.Value {
				p.take()
				return token
			}
		}
		unexpected(token.Value, values...)
	}
	unexpected(token.Type.String(), values...)
	panic("unreachable")
}

func (p *Parser) peek() Token {
	return p.tokens[p.position]
}

func (p *Parser) take() Token {
	token := p.tokens[p.position]

	if p.position != len(p.tokens)-1 {
		p.position++
	}

	return token
}

// Expect a token of TokenType t from the Parser's Tokens. If the type of the
// current token does not match t a panic will occur
func (p *Parser) expect(t TokenType) Token {
	token := p.peek()

	if token.Type != t {
		unexpected(token.Type.String(), t.String())
	}
	p.take()

	return token
}

func (p *Parser) optional(t TokenType) (Token, bool) {
	token := p.peek()

	if token.Type == t {
		p.take()
		return token, true
	}

	return InvalidToken, false
}

func unexpected(actual string, expected ...string) {
	panic(fmt.Errorf("Expected %s but found %s", strings.Join(expected, " or "), actual))
}
