package schema

import "fmt"

// Schema describes the structure and behavior of a GraphQL service
type Schema struct {
	enums      map[string]*Enum
	inputs     map[string]*Input
	interfaces map[string]*Interface
	objects    map[string]*Object
	scalars    map[string]*Scalar
	unions     map[string]*Union
}

// Declaration represents a declared Type in the GraphQL Schema
type Declaration interface {
	TypeKind() TypeKind
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

// TypeKind returns the TypeKind of Scalar
func (scalar Scalar) TypeKind() TypeKind {
	return SCALAR
}

// Enum describes an enum type within a Schema
type Enum struct {
	Name        string
	Values      []string
	Description string
}

// TypeKind returns the TypeKind of Enum
func (enum Enum) TypeKind() TypeKind {
	return ENUM
}

// Input describes an input type object within a Schema
type Input struct {
	Name        string
	Description string
	Fields      map[string]Field
}

// TypeKind returns the TypeKind of Input
func (input Input) TypeKind() TypeKind {
	return INPUT_OBJECT
}

// Interface describes an object type interface within a Schema
type Interface struct {
	Name        string
	Description string
	Fields      map[string]Field
}

// TypeKind returns the TypeKind of Interface
func (i Interface) TypeKind() TypeKind {
	return INTERFACE
}

// Union describes a union of Types within a Schema
type Union struct {
	Name        string
	Description string
	Types       []string
}

// TypeKind returns the TypeKind of Union
func (union Union) TypeKind() TypeKind {
	return UNION
}

// Object describes an Object Type defined within a Schema
type Object struct {
	Name        string
	Implements  []string
	Description string
	Fields      map[string]Field
}

// TypeKind returns the TypeKind of Object
func (object Object) TypeKind() TypeKind {
	return OBJECT
}

// Field describes a Field for a Type defined within a Schema
type Field struct {
	Name        string
	Description string
	Type        Type
	Arguments   map[string]Argument
}

// Type represents a Type in a Schema
type Type struct {
	Name    string
	NonNull bool
	List    bool
	SubType *Type
}

// Argument defines an argument for a Field defined within a Type
type Argument struct {
	Name    string
	Type    Type
	Default interface{}
}

///
// Schema Builder
///

// Builder builds a Schema from the Schema Declarations
type Builder struct {
	enums      []*Enum
	inputs     []*Input
	interfaces []*Interface
	objects    []*Object
	scalars    []*Scalar
	unions     []*Union
	errors     []error
}

// NewBuilder returns a new Builder
func NewBuilder() *Builder {
	return &Builder{}
}

func (builder *Builder) err(format string, s ...interface{}) {
	var err error

	if len(s) == 0 {
		err = fmt.Errorf("schema validation error: %s", format)
	} else {
		err = fmt.Errorf("schema validation error: %s", fmt.Sprintf(format, s))
	}
	builder.errors = append(builder.errors, err)
}

// Build builds and validates the Schema. If there are any errors a nil Schema
// will be returned along with the list of validation errors
func (builder *Builder) Build() (*Schema, []error) {
	var schema Schema

	// First, load all the Schema type maps with the declarations, and ensure
	// each one has only been defined once. Since each declaration must have a
	// name, any declaration without a name will be skipped and an error will be
	// noted
	schema.enums = make(map[string]*Enum, len(builder.enums))
	for _, enum := range builder.enums {
		if enum.Name == "" {
			builder.err("Enum declared without a Name")
			continue
		}

		if _, exists := schema.enums[enum.Name]; exists {
			builder.err("Enum %s declared multiple times", enum.Name)
			continue
		}
		schema.enums[enum.Name] = enum
	}

	schema.inputs = make(map[string]*Input, len(builder.inputs))
	for _, input := range builder.inputs {
		if builder.validateInputStructure(input) {
			if _, exists := schema.inputs[input.Name]; exists {
				builder.err("Input %s declared multiple times", input.Name)
				continue
			}
			schema.inputs[input.Name] = input
		}
	}

	schema.interfaces = make(map[string]*Interface, len(builder.interfaces))
	for _, i := range builder.interfaces {
		if i.Name == "" {
			builder.err("Interface declared without a Name")
			continue
		}

		if _, exists := schema.interfaces[i.Name]; exists {
			builder.err("Interface %s declared multiple times", i.Name)
			continue
		}
		schema.interfaces[i.Name] = i
	}

	schema.objects = make(map[string]*Object, len(builder.objects))
	for _, object := range builder.objects {
		if object.Name == "" {
			builder.err("Object declared without a Name")
			continue
		}

		if _, exists := schema.objects[object.Name]; exists {
			builder.err("Object %s declared multiple times", object.Name)
			continue
		}
		schema.objects[object.Name] = object
	}

	schema.scalars = make(map[string]*Scalar, len(builder.scalars))
	for _, scalar := range builder.scalars {
		if scalar.Name == "" {
			builder.err("Scalar declared without a Name")
			continue
		}

		if _, exists := schema.scalars[scalar.Name]; exists {
			builder.err("Scalar %s declared multiple times", scalar.Name)
			continue
		}
		schema.scalars[scalar.Name] = scalar
	}

	schema.unions = make(map[string]*Union, len(builder.unions))
	for _, union := range builder.unions {
		if union.Name == "" {
			builder.err("Union declared without a Name")
			continue
		}

		if _, exists := schema.unions[union.Name]; exists {
			builder.err("Union %s declared multiple times", union.Name)
			continue
		}
		schema.unions[union.Name] = union
	}
	// TODO a lot of copy paste above... can this be any better?

	if len(builder.errors) > 0 {
		return nil, builder.errors
	}
	return &schema, nil
}

// Enum adds a new Enum type declaration to the Schema
func (builder *Builder) Enum(enum Enum) *Builder {
	builder.enums = append(builder.enums, &enum)
	return builder
}

// Input adds a new Input type declaration to the Schema
func (builder *Builder) Input(input Input) *Builder {
	builder.inputs = append(builder.inputs, &input)
	return builder
}

// Interface adds a new Interface type declaration to the Schema
func (builder *Builder) Interface(i Interface) *Builder {
	builder.interfaces = append(builder.interfaces, &i)
	return builder
}

// Object adds a new Object type declaration to the Schema
func (builder *Builder) Object(object Object) *Builder {
	builder.objects = append(builder.objects, &object)
	return builder
}

// Scalar adds a new Scalar type declaration to the Schema
func (builder *Builder) Scalar(scalar Scalar) *Builder {
	builder.scalars = append(builder.scalars, &scalar)
	return builder
}

// Union adds a new Union type declaration to the Schema
func (builder *Builder) Union(union Union) *Builder {
	builder.unions = append(builder.unions, &union)
	return builder
}

// Declare a new schema type
func (builder *Builder) Declare(declaration Declaration) *Builder {
	switch v := declaration.(type) {
	case Enum:
		builder.Enum(v)
	case Input:
		builder.Input(v)
	case Interface:
		builder.Interface(v)
	case Object:
		builder.Object(v)
	case Scalar:
		builder.Scalar(v)
	case Union:
		builder.Union(v)
	default:
		// NOTE: unreachable only if type switch covers all Declaration types
		panic("unreachable")
	}

	return builder
}

func (builder *Builder) validateInputStructure(input *Input) (valid bool) {
	if input.Name == "" {
		builder.err("Input declared without a Name")
	} else if len(input.Fields) == 0 {
		builder.err("Input %s was declared without any fields defined", input.Name)
	} else {
		valid = true
	}
	return
}

func (builder *Builder) validateInputDeclaration(input *Input, schema *Schema) {
	if len(input.Fields) == 0 {
		builder.err("Input %s was declared without any fields defined", input.Name)
	}

}

// An Interface type must define one or more fields.
// The fields of an Interface type must have unique names within that Interface type; no two fields may share the same name.
func (builder *Builder) validateInterface(intrface *Interface, schema *Schema) {
	if len(intrface.Fields) == 0 {
		builder.err("Interface %s was declared without any fields defined", intrface.Name)
	}
}

// Object Validation Rules:
//
// An Object type must define one or more fields.
// The fields of an Object type must have unique names within that Object type; no two fields may share the same name.
// An object type must be a super‐set of all interfaces it implements:
// 		The object type must include a field of the same name for every field defined in an interface.
// 				The object field must be of a type which is equal to or a sub‐type of the interface field (covariant).
// 						An object field type is a valid sub‐type if it is equal to (the same type as) the interface field type.
// 						An object field type is a valid sub‐type if it is an Object type and the interface field type is either an Interface type or a Union type and the object field type is a possible type of the interface field type.
// 						An object field type is a valid sub‐type if it is a List type and the interface field type is also a List type and the list‐item type of the object field type is a valid sub‐type of the list‐item type of the interface field type.
// 						An object field type is a valid sub‐type if it is a Non‐Null variant of a valid sub‐type of the interface field type.
//				The object field must include an argument of the same name for every argument defined in the interface field.
// 						The object field argument must accept the same type (invariant) as the interface field argument.
// 				The object field may include additional arguments not defined in the interface field, but any additional argument must not be required.
func (builder *Builder) validateObject(object *Object, schema *Schema) {
	if len(object.Fields) == 0 {
		builder.err("Object %s was declared without any fields defined", object.Name)
	}

	// for _, interfaceName := object.Implements {
	// 	if intrface, exists := schema.interfaces[interfaceName]; exists {
	//
	// 	} else {
	// 		builder.err("Object %s declared to implement unknown Interface %s", object.Name, interfaceName)
	// 	}
	// }
}

//The member types of a Union type must all be Object base types; Scalar, Interface and Union types may not be member types of a Union. Similarly, wrapping types may not be member types of a Union.
// A Union type must define one or more member types.
func (builder *Builder) validateUnion(union *Union, schema *Schema) {
	if len(union.Types) == 0 {
		builder.err("Union %s was declared without any member types defined")
	}

	// All member types must be Object types
	for _, member := range union.Types {
		if _, exists := schema.objects[member]; !exists {
			builder.err("Union %s was declared with the member %s which is not an Object type", union.Name, member)
		}
	}
}

///
// Schema Builder Helpers (mostly for readability purposes)
///

// Interfaces builds a list of Interfaces
func Interfaces(interfaces ...string) []string {
	return interfaces
}

// Objects builds a list of Objects
func Objects(objects ...Object) []Object {
	return objects
}

// Arguments builds a map from Argument names to the Argument itself
func Arguments(arguments ...Argument) map[string]Argument {
	argumentsMap := make(map[string]Argument)

	for _, argument := range arguments {
		argumentsMap[argument.Name] = argument
	}
	return argumentsMap
}

// Fields builds a map from Field names to the Field itself
func Fields(fields ...Field) map[string]Field {
	fieldsMap := make(map[string]Field)

	for _, field := range fields {
		fieldsMap[field.Name] = field
	}
	return fieldsMap
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
