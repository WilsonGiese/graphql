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

type Directives struct {
}

type SelectionSet struct {
}
