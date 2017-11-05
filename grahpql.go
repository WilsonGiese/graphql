package graphql

type Document struct {
}

type Operation struct {
	operationType string
	operationName string
	variables     []Variable
	directives    Directives
	selectionSet  SelectionSet
}

type Fragment struct {
}

type Variable struct {
	Name    string
	Type    Type
	Default Value
}

type ValueType int

const (
	NullType ValueType = iota
	VariableType
	IntegerType
	FloatType
	StringType
	BooleanType
	EnumType
	ListType
	ObjectType
)

type Value struct {
	Type  ValueType
	Value interface{}
}

type Type struct {
	Type    ValueType
	SubType *Type
	NonNull bool
}

// {Type: List, SubType: {Type: String, SubType: nil}}

type Object struct {
}

type Directives struct {
}

type SelectionSet struct {
}
