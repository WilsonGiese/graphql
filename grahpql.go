package graphql

import "errors"

type Document struct {
	Operations []Operation
	Fragments  []Fragment
}

func (document *Document) GetFragment(name string) (Fragment, error) {
	for _, fragment := range document.Fragments {
		if fragment.Name == name {
			return fragment, nil
		}
	}
	return Fragment{}, errors.New("No fragment found with that name")
}

type Operation struct {
	Type                string
	Name                string
	VariableDefinitions []VariableDefinition
	Directives          []Directive
	SelectionSet        SelectionSet
}

type VariableDefinition struct {
	Name    string
	Type    Type
	Default Value
}

type Value struct {
	Type  Type
	Value interface{}
}

type Type struct {
	Type    string
	NonNull bool
	List    bool
	SubType *Type
}

type Object struct {
}

type Directive struct {
	Name      string
	Arguments map[string]Value
}

type SelectionSet struct {
	Fields          []Field
	InlineFragments []InlineFragment
	FragmentSpreads []FragmentSpread
}

// IsEmpty returns true if the SelectionSet contains no field selections
// or fragments; false otherwise
func (selectionSet SelectionSet) IsEmpty() bool {
	return len(selectionSet.Fields) == 0 &&
		len(selectionSet.FragmentSpreads) == 0 &&
		len(selectionSet.InlineFragments) == 0
}

type Fragment struct {
	Name         string
	Type         string
	Directives   []Directive
	SelectionSet SelectionSet
}

type InlineFragment struct {
	Type         string
	Directives   []Directive
	SelectionSet SelectionSet
}

type FragmentSpread struct {
	Name       string
	Directives []Directive
}

type Field struct {
	Name         string
	Alias        string
	Arguments    map[string]Value
	Directives   []Directive
	SelectionSet SelectionSet
}
