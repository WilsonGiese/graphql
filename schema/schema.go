package schema

import (
	"errors"
	"fmt"
)

// Schema describes the structure and behavior of a GraphQL service
type Schema struct {
	enums      map[string]Enum
	inputs     map[string]Input
	interfaces map[string]Interface
	objects    map[string]Object
	scalars    map[string]Scalar
	unions     map[string]Union
}

func newSchema() *Schema {
	return &Schema{
		enums:      make(map[string]Enum),
		inputs:     make(map[string]Input),
		interfaces: make(map[string]Interface),
		objects:    make(map[string]Object),
		scalars:    make(map[string]Scalar),
		unions:     make(map[string]Union),
	}
}

// Declaration represents a declared Type in the GraphQL Schema
type Declaration interface {
	name() string
	typeKind() TypeKind
}

// TypeKind represents the type of a Schema declaration
type TypeKind int

// TypeKinds represents the complete set of types a Schema declaration can be
const (
	SCALAR TypeKind = iota
	OBJECT
	INTERFACE
	UNION
	ENUM
	INPUT_OBJECT
	LIST
)

///
// Schema Types
///

// Scalar describes a scalar type within a Schema
type Scalar struct {
	Name        string
	Description string
}

func (scalar Scalar) name() string {
	return scalar.Name
}

// TypeKind returns the TypeKind of Scalar
func (scalar Scalar) typeKind() TypeKind {
	return SCALAR
}

// Enum describes an enum type within a Schema
type Enum struct {
	Name        string
	Values      []string
	Description string
}

func (enum Enum) name() string {
	return enum.Name
}

// TypeKind returns the TypeKind of Enum
func (enum Enum) typeKind() TypeKind {
	return ENUM
}

func (enum Enum) String() string {
	return fmt.Sprintf("Enum(%s)", enum.Name)
}

// Input describes an input type object within a Schema
type Input struct {
	Name        string
	Description string
	Fields      map[string]Field
}

func (input Input) name() string {
	return input.Name
}

// TypeKind returns the TypeKind of Input
func (input Input) typeKind() TypeKind {
	return INPUT_OBJECT
}

func (input Input) String() string {
	return fmt.Sprintf("Input(%s)", input.Name)
}

// Interface describes an object type interface within a Schema
type Interface struct {
	Name        string
	Description string
	Fields      map[string]Field
}

func (intrface Interface) name() string {
	return intrface.Name
}

// TypeKind returns the TypeKind of Interface
func (intrface Interface) typeKind() TypeKind {
	return INTERFACE
}

func (intrface Interface) String() string {
	return fmt.Sprintf("Interface(%s)", intrface.Name)
}

// Union describes a union of Types within a Schema
type Union struct {
	Name        string
	Description string
	Types       []string
}

func (union Union) name() string {
	return union.Name
}

// TypeKind returns the TypeKind of Union
func (union Union) typeKind() TypeKind {
	return UNION
}

func (union Union) String() string {
	return fmt.Sprintf("Union(%s)", union.Name)
}

// Object describes an Object Type defined within a Schema
type Object struct {
	Name        string
	Implements  []string
	Description string
	Fields      map[string]Field
}

func (object Object) name() string {
	return object.Name
}

// TypeKind returns the TypeKind of Object
func (object Object) typeKind() TypeKind {
	return OBJECT
}

func (object Object) String() string {
	return fmt.Sprintf("Object(%s)", object.Name)
}

// Field describes a Field for a Type defined within a Schema
type Field struct {
	Name        string
	Description string
	Type        Type
	Arguments   map[string]Argument
}

func (field Field) String() string {
	return fmt.Sprintf("Field(%s)", field.Name)
}

// Type represents a Type in a Schema
type Type struct {
	Name    string
	NonNull bool
	List    bool
	SubType *Type
}

func (t Type) String() string {
	if t.List {
		if t.SubType == nil {
			return "[nil]"
		}
		if t.NonNull {
			return fmt.Sprintf("[%s]!", t.SubType)
		}
		return fmt.Sprintf("[%s]", t.SubType)
	}

	if t.NonNull {
		return t.Name + "!"
	}
	return t.Name
}

// Argument defines an argument for a Field defined within a Type
type Argument struct {
	Name    string
	Type    Type
	Default interface{}
}

func (argument Argument) String() string {
	return fmt.Sprintf("Argument(%s)", argument.Name)
}

var errTypeNotFound = errors.New("type not found")

func (schema *Schema) getEnum(name string) (Enum, error) {
	if enum, exists := schema.enums[name]; exists {
		return enum, nil
	}
	return Enum{}, errTypeNotFound
}

func (schema *Schema) getInput(name string) (Input, error) {
	if input, exists := schema.inputs[name]; exists {
		return input, nil
	}
	return Input{}, errTypeNotFound
}

func (schema *Schema) getInterface(name string) (Interface, error) {
	if intrface, exists := schema.interfaces[name]; exists {
		return intrface, nil
	}
	return Interface{}, errTypeNotFound
}

func (schema *Schema) getObject(name string) (Object, error) {
	if object, exists := schema.objects[name]; exists {
		return object, nil
	}
	return Object{}, errTypeNotFound
}

func (schema *Schema) getScalar(name string) (Scalar, error) {
	if scalar, exists := schema.scalars[name]; exists {
		return scalar, nil
	}
	return Scalar{}, errTypeNotFound
}

func (schema *Schema) getUnion(name string) (Union, error) {
	if union, exists := schema.unions[name]; exists {
		return union, nil
	}
	return Union{}, errTypeNotFound
}

// getDeclaration returns the Declaration for a Type. If the type is a list type
// it will return the unwrapped type; e.g. [[String]] -> String
func (schema *Schema) getDeclaration(t Type) Declaration {
	if t.List {
		return schema.getDeclaration(*t.SubType)
	}

	if scalar, exists := schema.scalars[t.Name]; exists {
		return scalar
	}

	if enum, exists := schema.enums[t.Name]; exists {
		return enum
	}

	if object, exists := schema.objects[t.Name]; exists {
		return object
	}

	if intrface, exists := schema.interfaces[t.Name]; exists {
		return intrface
	}

	if union, exists := schema.unions[t.Name]; exists {
		return union
	}

	if input, exists := schema.inputs[t.Name]; exists {
		return input
	}

	return nil
}

///
// Type helpers
///

// DescribeType returns a TypeSchema for the given type which can be null
func DescribeType(name string) (t Type) {
	t.Name = name
	return
}

// DescribeNonNullType returns a non-null TypeSchema for the given type
func DescribeNonNullType(name string) (t Type) {
	t.Name = name
	t.NonNull = true
	return
}

// DescribeListType returns a TypeSchema for the given subType which can be null
func DescribeListType(subType Type) (t Type) {
	t.List = true
	t.SubType = &subType
	return
}

// DescribeNonNullListType returns a non-null TypeSchema for the given subType
func DescribeNonNullListType(subType Type) (t Type) {
	t.List = true
	t.NonNull = true
	t.SubType = &subType
	return
}

///
// Common Pre-defined Types (Int, Float, String, Boolean, ID)
///

// StringType is a TypeSchema for a String which can be null
var StringType = DescribeType("String")

// NonNullStringType is a TypeSchema for String a string which cannot be null
var NonNullStringType = DescribeNonNullType("String")

// IntType is a TypeSchema for an Int which can be null
var IntType = DescribeType("Int")

// NonNullIntType is a TypeSchema for an Int which cannot be null
var NonNullIntType = DescribeNonNullType("Int")

// FloatType is a TypeSchema for a Float which can be null
var FloatType = DescribeType("Float")

// NonNullFloatType is a TypeSchema for a Float which cannot be null
var NonNullFloatType = DescribeNonNullType("Float")

// BooleanType is a TypeSchema for a Boolean which can be null
var BooleanType = DescribeType("Boolean")

// NonNullBooleanType is a TypeSchema for a Boolean which cannot be null
var NonNullBooleanType = DescribeNonNullType("Boolean")

// IDType is a TypeSchema for an ID which can be null
var IDType = DescribeType("ID")

// NonNullIDType is a TypeSchema for an ID which cannot be null
var NonNullIDType = DescribeNonNullType("ID")
