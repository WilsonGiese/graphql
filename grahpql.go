package graphql

type Document struct {
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

type ValueType int

const (
	ObjectType ValueType = iota
	IntegerType
	FloatType
	StringType
	BooleanType
	EnumType
	ListType
)

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

// {Type: List, SubType: {Type: String, SubType: nil}}

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
	Type       string
	Directives []Directive
}

type Field struct {
	Name         string
	Alias        string
	Arguments    map[string]Value
	Directives   []Directive
	SelectionSet SelectionSet
}
