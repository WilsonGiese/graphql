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

func (p *Parser) parse() (document Document, err error) {

	// Defer function to recover from a panic during parsing, and finally return
	// the panic back as a proper error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	document = p.parseDocument()
	return
}

func (p *Parser) parseDocument() (document Document) {
	// A document starts with an OpenBrace if it is using query short-hand syntax,
	// otherwise a document is a list of Operations with their type stated
	if p.peek().Type == OpenBrace {
		document.Operations = append(document.Operations, Operation{SelectionSet: p.parseSelectionSet()})
	} else {
		for {
			token := p.peek()
			if token.Type == EOF {
				break
			}
			defintionType := p.accept(Name, "query", "mutation", "fragment").Value

			if defintionType == "fragment" {
				document.Fragments = append(document.Fragments, p.parseFragment())
			} else {
				document.Operations = append(document.Operations, p.parseOperation(defintionType))
			}
		}
	}
	p.expect(EOF)

	return Document{}
}

// OperationDefinition
// OperationType Name(opt) VariableDefinitions(opt) Directives(opt) SelectionSet
// SelectionSet
func (p *Parser) parseOperation(operationType string) (operation Operation) {
	operation.Type = operationType

	if name, containsName := p.optional(Name); containsName {
		operation.Name = name.Value
	}

	operation.VariableDefinitions = p.parseVariableDefinitions()
	operation.Directives = p.parseDirectives()
	operation.SelectionSet = p.parseSelectionSet()
	return
}

func (p *Parser) parseFragment() (fragment Fragment) {
	fragment.Name = p.expect(Name).Value
	fragment.Type = p.parseTypeCondition()
	fragment.Directives = p.parseDirectives()
	fragment.SelectionSet = p.parseSelectionSet()
	return
}

func (p *Parser) parseTypeCondition() string {
	if token := p.peek(); token.Type == Name {
		if token.Value != "on" {
			unexpected(token.Value, "on")
		}
		p.take()
	} else {
		unexpected(token.Type.String(), "on")
	}
	return p.expect(Name).Value
}

func (p *Parser) parseVariableDefinitions() (varDefs []VariableDefinition) {
	p.expect(OpenParen)
	for {
		token := p.peek()

		if token.Type != Dollar {
			break
		}

		varDefs = append(varDefs, p.parseVariableDefinition())
	}
	p.expect(ClosedParen)
	return
}

func (p *Parser) parseVariableDefinition() (varDef VariableDefinition) {
	p.expect(Dollar)
	varDef.Name = p.expect(Name).Value
	p.expect(Colon)
	varDef.Type = p.parseType()

	if _, defaultGiven := p.optional(Equals); defaultGiven {
		varDef.Default = p.parseValue()
	}
	return
}

func (p *Parser) parseType() (t Type) {
	if p.peek().Type == OpenBracket {
		p.expect(OpenBracket)
		subType := p.parseType()
		p.expect(ClosedBracket)

		t.List = true
		t.SubType = &subType
	} else {
		t.Type = p.expect(Name).Value
	}

	_, nonNull := p.optional(Exclamation)
	t.NonNull = nonNull

	return
}

func (p *Parser) parseValue() (v Value) {
	token := p.peek()
	switch p.peek().Type {
	case Name:
		fallthrough
	case String:
		fallthrough
	case Integer:
		fallthrough
	case Float:
		p.take()
		v.Value = token.Value
	case OpenBracket:
		v.Value = p.parseListValue()
	case OpenBrace:
		v.Value = p.parseObjectValue()
	default:
		unexpected(token.Type.String(), "Value")
	}
	return
}

func (p *Parser) parseListValue() (values []Value) {
	p.expect(OpenBracket)
	for {
		if p.peek().Type == ClosedBracket {
			break
		}

		values = append(values, p.parseValue())
	}
	p.expect(ClosedBracket)
	return
}

func (p *Parser) parseObjectValue() (object map[string]Value) {
	p.expect(OpenBrace)
	for {
		if p.peek().Type == ClosedBrace {
			break
		}

		name := p.expect(Name)
		p.expect(Colon)
		value := p.parseValue()

		if _, exists := object[name.Value]; exists {
			invalid("duplicate field name in object value")
		}
		object[name.Value] = value
	}
	p.expect(ClosedBrace)
	return
}

func (p *Parser) parseDirectives() (directives []Directive) {
	for {
		if p.peek().Type != At {
			break
		}

		directives = append(directives, p.parseDirective())
	}
	return
}

func (p *Parser) parseDirective() (directive Directive) {
	p.expect(At)
	name := p.expect(Name).Value
	arguments := p.parseArguments()

	directive.Name = name
	directive.Arguments = arguments
	return
}

func (p *Parser) parseArguments() (arguments map[string]Value) {
	p.expect(OpenParen)
	for {
		if p.peek().Type == Name {
			break
		}

		name := p.expect(Name)
		p.expect(Colon)
		value := p.parseValue()

		if _, exists := arguments[name.Value]; exists {
			invalid("duplicate argument in arguments list")
		}
		arguments[name.Value] = value
	}
	p.expect(ClosedParen)
	return
}

func (p *Parser) parseSelectionSet() (selectionSet SelectionSet) {
	p.expect(OpenBrace)
	for {
		token := p.peek()

		if token.Type != Name && token.Type != Spread {
			break
		}

		// Spread indicates a Fragment Spread or an Inline Fragment
		// Fragment Spread must start with a Name that is not "on".
		// Inline Fragment optionally starts with a type condition (on Name),
		// followed optionally by directives, followed finally by a selection set
		if token.Type == Spread {
			lookahead := p.lookahead(1)
			if lookahead.Type == Name {
				if lookahead.Value == "on" {
					selectionSet.InlineFragments = append(selectionSet.InlineFragments, p.parseInlineFragment())
				} else {
					selectionSet.FragmentSpreads = append(selectionSet.FragmentSpreads, p.parseFragmentSpread())
				}
			} else if lookahead.Type == At || lookahead.Type == OpenBrace {
				selectionSet.InlineFragments = append(selectionSet.InlineFragments, p.parseInlineFragment())
			} else {
				unexpected(lookahead.Type.String(), "fragment spread or inline fragment")
			}
		} else {
			selectionSet.Fields = append(selectionSet.Fields, p.parseField())
		}
	}
	p.expect(ClosedBrace)
	return
}

func (p *Parser) parseFragmentSpread() (fragmentSpread FragmentSpread) {
	p.expect(Spread)
	fragmentSpread.Type = p.expect(Name).Value
	fragmentSpread.Directives = p.parseDirectives()

	return
}

func (p *Parser) parseInlineFragment() (inlineFragment InlineFragment) {
	p.expect(Spread)
	inlineFragment.Type = p.parseTypeCondition()
	inlineFragment.Directives = p.parseDirectives()
	inlineFragment.SelectionSet = p.parseSelectionSet()
	return
}

func (p *Parser) parseField() (field Field) {
	field.Name = p.expect(Name).Value

	if p.peek().Type == Colon {
		field.Alias = field.Name
		field.Name = p.expect(Name).Value
	}

	if p.peek().Type == OpenParen {
		field.Arguments = p.parseArguments()
	}

	// Directives
	if p.peek().Type == At {
		field.Directives = p.parseDirectives()
	}

	if p.peek().Type == OpenBrace {
		field.SelectionSet = p.parseSelectionSet()
	}
	return
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

func (p *Parser) lookahead(distance int) Token {
	if p.position+distance >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}
	return p.tokens[p.position+distance]
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

func invalid(message string) {
	panic(fmt.Errorf("invalid: %s", message))
}
