package schema

import (
	"fmt"
	"reflect"
	"regexp"
)

var validNameMatcher *regexp.Regexp

func init() {
	validNameMatcher = regexp.MustCompile(`^[_a-zA-Z0-9]+$`)
}

// Builder builds a Schema from provided Declarations
type Builder struct {
	schema            *Schema
	declaredTypeNames map[string]interface{}
}

// NewSchema returns a new Schema Builder
func NewSchema() *Builder {
	builder := Builder{
		schema:            newSchema(),
		declaredTypeNames: make(map[string]interface{}),
	}

	// Default Schema
	builder.Declare(Object{
		Name: "__Schema",
		Fields: Fields(
			Field{
				Name:        "types",
				Description: "All types that are a part of this Schema",
				Type:        DescribeNonNullListType(DescribeNonNullType("__Type")),
			},
			Field{
				Name:        "queryType",
				Description: "Root query type for this Schema",
				Type:        DescribeNonNullType("__Type"),
			},
			Field{
				Name:        "mutationType",
				Description: "Root mutation type for this Schema",
				Type:        DescribeType("__Type"),
			},
			Field{
				Name:        "directives",
				Description: "All directives that are a part of this Schema",
				Type:        DescribeNonNullListType(DescribeNonNullType("__Directive")),
			},
		),
	}).Declare(Object{
		Name: "__Type",
		Fields: Fields(
			Field{
				Name: "kind",
				Type: DescribeNonNullType("__TypeKind"),
			},
			Field{
				Name: "name",
				Type: NonNullStringType,
			},
			Field{
				Name: "description",
				Type: StringType,
			},
			Field{
				Name: "fields",
				Type: DescribeNonNullType("__Field"),
				Arguments: Arguments(
					Argument{
						Name:    "includeDeprecated",
						Type:    BooleanType,
						Default: false,
					},
				),
			},
			Field{
				Name: "interfaces",
				Type: DescribeListType(DescribeNonNullType("__Type")),
			},
			Field{
				Name: "possibleTypes",
				Type: DescribeListType(DescribeNonNullType("__Type")),
			},
			Field{
				Name: "enumValues",
				Type: DescribeNonNullType("__EnumValue"),
				Arguments: Arguments(
					Argument{
						Name:    "includeDeprecated",
						Type:    BooleanType,
						Default: false,
					},
				),
			},
			Field{
				Name: "inputFields",
				Type: DescribeListType(DescribeNonNullType("__InputValue")),
			},
			Field{
				Name: "ofType",
				Type: DescribeType("__Type"),
			},
		),
	}).Declare(Object{
		Name: "__Field",
		Fields: Fields(
			Field{
				Name: "name",
				Type: NonNullStringType,
			},
			Field{
				Name: "description",
				Type: StringType,
			},
			Field{
				Name: "args",
				Type: DescribeNonNullListType(DescribeNonNullType("__InputValue")),
			},
			Field{
				Name: "type",
				Type: DescribeNonNullType("__Type"),
			},
			Field{
				Name: "isDeprecated",
				Type: NonNullBooleanType,
			},
			Field{
				Name: "deprecationReason",
				Type: StringType,
			},
		),
	}).Declare(Object{
		Name: "__InputValue",
		Fields: Fields(
			Field{
				Name: "name",
				Type: NonNullStringType,
			},
			Field{
				Name: "description",
				Type: StringType,
			},
			Field{
				Name: "type",
				Type: DescribeNonNullType("__Type"),
			},
			Field{
				Name: "defaultValue",
				Type: StringType,
			},
		),
	}).Declare(Object{
		Name: "__EnumValue",
		Fields: Fields(
			Field{
				Name: "name",
				Type: NonNullStringType,
			},
			Field{
				Name: "description",
				Type: StringType,
			},
			Field{
				Name: "isDeprecated",
				Type: NonNullBooleanType,
			},
			Field{
				Name: "deprecationReason",
				Type: StringType,
			},
		),
	}).Declare(Object{
		Name: "__Directive",
		Fields: Fields(
			Field{
				Name: "name",
				Type: NonNullStringType,
			},
			Field{
				Name: "description",
				Type: StringType,
			},
			Field{
				Name: "location",
				Type: DescribeNonNullListType(DescribeNonNullType("__DirectiveLocation")),
			},
			Field{
				Name: "args",
				Type: DescribeNonNullListType(DescribeNonNullType("__InputValue")),
			},
		),
	}).Declare(Enum{
		Name: "__TypeKind",
		Values: Values(
			"SCALAR",
			"OBJECT",
			"INTERFACE",
			"UNION",
			"ENUM",
			"INPUT_OBJECT",
			"LIST",
			"NON_NULL",
		),
	}).Declare(Enum{
		Name: "__DirectiveLocation",
		Values: Values(
			"QUERY",
			"MUTATION",
			"FIELD",
			"FRAGMENT_DEFINITION",
			"FRAGMENT_SPREAD",
			"INLINE_FRAGMENT",
		),
	}).Declare(Scalar{
		Name:        "Int",
		Description: "The Int scalar type represents a signed 32‐bit numeric non‐fractional value. Response formats that support a 32‐bit integer or a number type should use that type to represent this scalar.",
	}).Declare(Scalar{
		Name:        "Float",
		Description: "The Float scalar type represents signed double‐precision fractional values as specified by IEEE 754. Response formats that support an appropriate double‐precision number type should use that type to represent this scalar.",
	}).Declare(Scalar{
		Name:        "String",
		Description: "The String scalar type represents textual data, represented as UTF‐8 character sequences. The String type is most often used by GraphQL to represent free‐form human‐readable text. All response formats must support string representations, and that representation must be used here.",
	}).Declare(Scalar{
		Name:        "Boolean",
		Description: "The Boolean scalar type represents true or false. Response formats should use a built‐in boolean type if supported; otherwise, they should use their representation of the integers 1 and 0.",
	}).Declare(Scalar{
		Name:        "ID",
		Description: "The ID scalar type represents a unique identifier, often used to refetch an object or as the key for a cache. The ID type is serialized in the same way as a String; however, it is not intended to be human‐readable. While it is often numeric, it should always serialize as a String.",
	})

	return &builder
}

func (builder *Builder) err(format string, s ...interface{}) {
	var err error

	if len(s) == 0 {
		err = NewValidationError(format)
	} else {
		err = NewValidationError(fmt.Sprintf(format, s...))
	}
	panic(err)
}

// Build builds and validates the Schema. If there are any validation issues
// Build will panic with a schema validation error describing the problem
func (builder *Builder) Build() *Schema {
	for _, scalar := range builder.schema.scalars {
		builder.validateScalar(scalar)
	}
	for _, enum := range builder.schema.enums {
		builder.validateEnum(enum)
	}
	for _, intrface := range builder.schema.interfaces {
		builder.validateInterface(intrface)
	}
	for _, object := range builder.schema.objects {
		builder.validateObject(object)
	}
	for _, union := range builder.schema.unions {
		builder.validateUnion(union)
	}
	for _, input := range builder.schema.inputs {
		builder.validateInput(input)
	}
	return builder.schema
}

// Enum adds a new Enum type declaration to the Schema
func (builder *Builder) Enum(enum Enum) *Builder {
	if err := builder.declareTypeName(enum); err != nil {
		builder.err("Enum %s", err)
	}
	builder.schema.enums[enum.Name] = enum
	return builder
}

// Input adds a new Input type declaration to the Schema
func (builder *Builder) Input(input Input) *Builder {
	if err := builder.declareTypeName(input); err != nil {
		builder.err("Input %s", err)
	}
	builder.schema.inputs[input.Name] = input
	return builder
}

// Interface adds a new Interface type declaration to the Schema
func (builder *Builder) Interface(intrface Interface) *Builder {
	if err := builder.declareTypeName(intrface); err != nil {
		builder.err("Interface %s", err)
	}
	builder.schema.interfaces[intrface.Name] = intrface
	return builder
}

// Object adds a new Object type declaration to the Schema. A panic will occur
// if the Object is declared with an invalid Name, or a Type with that Name has
// already been declared
func (builder *Builder) Object(object Object) *Builder {
	if err := builder.declareTypeName(object); err != nil {
		builder.err("Object %s", err)
	}
	builder.schema.objects[object.Name] = object
	return builder
}

// Scalar adds a new Scalar type declaration to the Schema
func (builder *Builder) Scalar(scalar Scalar) *Builder {
	if err := builder.declareTypeName(scalar); err != nil {
		builder.err("Scalar %s", err)
	}
	builder.schema.scalars[scalar.Name] = scalar
	return builder
}

// Union adds a new Union type declaration to the Schema
func (builder *Builder) Union(union Union) *Builder {
	if err := builder.declareTypeName(union); err != nil {
		builder.err("Union %s", err)
	}
	builder.schema.unions[union.Name] = union
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

func (builder *Builder) declareTypeName(declaration Declaration) error {
	if err := builder.validateName(declaration.name()); err != nil {
		return err
	}

	if _, exists := builder.declaredTypeNames[declaration.name()]; exists {
		return fmt.Errorf("declared with name '%s' but another type with that name has already been declared", declaration.name())
	}
	builder.declaredTypeNames[declaration.name()] = struct{}{}
	return nil
}

///
// Validation Functions - ensure schema declarations are valid
///

// Interface 'Pet': Field 'name': declared with unknown Type 'String'

func (builder *Builder) validateEnum(enum Enum) {
	if len(enum.Values) == 0 {
		builder.err("%s delcared without any values defined", enum)
	}

	if duplicate := findFirstDuplicate(enum.Values); duplicate != nil {
		builder.err("%s declared duplicate value %s", enum, *duplicate)
	}
}

// http://facebook.github.io/graphql/October2016/#sec-Input-Object-type-validation
func (builder *Builder) validateInput(input Input) {
	if len(input.Fields) == 0 {
		builder.err("%s declared without any Fields defined", input)
	}

	for _, field := range input.Fields {
		if err := builder.validateInputField(field); err != nil {
			builder.err("%s %s", input, err)
		}
	}
}

func (builder *Builder) validateInputField(field Field) error {
	if err := builder.validateName(field.Name); err != nil {
		return fmt.Errorf("Field %s", err)
	}

	if err := builder.validateType(field.Type); err != nil {
		return fmt.Errorf("%s %s", field, err)
	}

	// Input field types can only be input, scalar, or enum
	switch builder.schema.getDeclaration(field.Type).typeKind() {
	case INPUT_OBJECT:
	case SCALAR:
	case ENUM:
	default:
		return fmt.Errorf("%s declared with invalid Type '%s'. An Input Field type must be Input, Scalar, or Enum", field, field.Type)
	}

	// Input fields cannot be declared with arguments
	if len(field.Arguments) > 0 {
		return fmt.Errorf("%s declared with arguments. Input fields must be declared without arguments", field)
	}
	return nil
}

// http://facebook.github.io/graphql/October2016/#sec-Interface-type-validation
func (builder *Builder) validateInterface(intrface Interface) {
	if len(intrface.Fields) == 0 {
		builder.err("%s declared without any Fields defined", intrface)
	}

	for _, field := range intrface.Fields {
		if err := builder.validateField(field); err != nil {
			builder.err("%s %s", intrface, err)
		}
	}
}

// http://facebook.github.io/graphql/October2016/#sec-Object-type-validation
func (builder *Builder) validateObject(object Object) {

	if len(object.Fields) == 0 {
		builder.err("%s declared without any Fields defined", object)
	}

	for _, field := range object.Fields {
		if err := builder.validateField(field); err != nil {
			builder.err("%s %s", object, err)
		}
	}

	for _, interfaceName := range object.Implements {
		if intrface, err := builder.schema.getInterface(interfaceName); err == nil {
			for _, interfaceField := range intrface.Fields {
				if objectField, exists := object.Fields[interfaceField.Name]; exists {
					if err := builder.validateFieldImplementsInterface(objectField, interfaceField, intrface); err != nil {
						builder.err("%s %s", object, err)
					}
				} else {
					builder.err("%s declared without %s required by %s", object, interfaceField, intrface)
				}
			}
		} else {
			builder.err("%s declared implementing unknown Interface '%s'", object, interfaceName)
		}
	}
}

func (builder *Builder) validateScalar(scalar Scalar) {
	// NOTE: Nothing to validate for now
}

// http://facebook.github.io/graphql/October2016/#sec-Union-type-validation
func (builder *Builder) validateUnion(union Union) {
	if len(union.Types) == 0 {
		builder.err("%s declared without any member types defined", union)
	}

	for _, unionTypeName := range union.Types {
		if declaration := builder.schema.getDeclaration(DescribeType(unionTypeName)); declaration == nil {
			builder.err("%s declared with unknown type %s", union, unionTypeName)
		} else if declaration.typeKind() != OBJECT {
			builder.err("%s declared with member type %s. Union members must be Objects", union, unionTypeName)
		}
	}

	if duplicate := findFirstDuplicate(union.Types); duplicate != nil {
		builder.err("%s declared duplicate type %s", union, *duplicate)
	}
}

// Vaidate a Object or Interface field. Use validateInputField for Input types
func (builder *Builder) validateField(field Field) error {
	if err := builder.validateName(field.Name); err != nil {
		return fmt.Errorf("Field %s", err)
	}

	if err := builder.validateType(field.Type); err != nil {
		return fmt.Errorf("%s %s", field, err)
	}

	if builder.schema.getDeclaration(field.Type).typeKind() == INPUT_OBJECT {
		return fmt.Errorf("%s declared with Input type '%s'", field, field.Type)
	}

	for _, argument := range field.Arguments {
		if err := builder.validateArgument(argument); err != nil {
			return fmt.Errorf("%s %s", field, err)
		}
	}
	return nil
}

func (builder *Builder) validateFieldImplementsInterface(objectField, interfaceField Field, intrface Interface) error {
	if builder.covariantTypeCheck(objectField.Type, interfaceField.Type) {

		// Validate all interface Field Arguments are implemented
		for _, interfaceArgument := range interfaceField.Arguments {
			if objectArgument, exists := objectField.Arguments[interfaceArgument.Name]; exists {
				if err := builder.validateArgumentImplementsInterface(objectArgument, interfaceArgument, intrface); err != nil {
					return fmt.Errorf("%s %s", objectField, err)
				}
			} else {
				return fmt.Errorf("%s declared without %s required by %s", objectField, interfaceArgument, intrface)
			}
		}

		// Validate all aditional Field Arguments are not required
		// TODO this is wrong, idiot! If THIS interface doesn't require this field, ANOTHER might.
		for _, objectArgument := range objectField.Arguments {
			if _, exists := interfaceField.Arguments[objectArgument.Name]; !exists {
				if objectArgument.Type.NonNull {
					return fmt.Errorf("%s declared an additional %s with a non-null type. Since %s is required by %s any additional Arguments must not be required", objectField, objectArgument, interfaceField, intrface)
				}
			}
		}
	} else {
		return fmt.Errorf("%s declared with type '%s' but %s requires the type '%s' or a valid sub-type", objectField, objectField.Type, intrface, interfaceField.Type)
	}
	return nil
}

// - The object field must include an argument of the same name for every
//   argument defined in the interface field.
// 		- The object field argument must accept the same type (invariant) as the
//      interface field argument.
// - The object field may include additional arguments not defined in the
//   interface field, but any additional argument must not be required.
func (builder *Builder) validateArgumentImplementsInterface(objectArg, interfaceArg Argument, intrface Interface) error {
	if !typeCheck(objectArg.Type, interfaceArg.Type) {
		return fmt.Errorf("%s declared with type '%s' but %s requires type '%s'", objectArg, objectArg.Type, intrface, interfaceArg.Type)
	}
	return nil
}

func (builder *Builder) validateArgument(argument Argument) error {
	if err := builder.validateName(argument.Name); err != nil {
		return fmt.Errorf("Argument: %s", err)
	}

	if err := builder.validateType(argument.Type); err != nil {
		return fmt.Errorf("%s %s", argument, err)
	}

	// Input field types can only be input, scalar, or enum
	switch builder.schema.getDeclaration(argument.Type).typeKind() {
	case INPUT_OBJECT:
	case SCALAR:
	case ENUM:
	default:
		return fmt.Errorf("%s declared with invalid type '%s'. An Argument Type must be Input, Scalar, or Enum", argument, argument.Type)
	}

	if argument.Type.NonNull && argument.Default != nil {
		return fmt.Errorf("%s declared with a default value, but its type is non-null", argument)
	}
	return nil
	// TODO validate default value?
}

func (builder *Builder) validateName(name string) error {
	if name == "" {
		return fmt.Errorf("declared without Name defined")
	}

	if !validNameMatcher.MatchString(name) {
		return fmt.Errorf("declared with an invalid Name '%s'. A Name must only consist of ASCII letters, numbers, and underscores", name)
	}
	return nil
}

func (builder *Builder) validateType(t Type) error {
	// If t is a list type, pull out the underlying base type
	actualType := t
	for {
		if !actualType.List {
			break
		}

		if t.SubType == nil {
			return fmt.Errorf("type %s has nil sub-type", t)
		}
		actualType = *t.SubType
	}

	if declaration := builder.schema.getDeclaration(actualType); declaration == nil {
		return fmt.Errorf("declared with unknown type '%s'", actualType.Name)
	}
	return nil
}

// typeCheck checks if Types t1 and t2 are equal to eachother
func typeCheck(t1, t2 Type) bool {
	return reflect.DeepEqual(t1, t2)
}

// covariantTypeCheck checks if type t1 is equal to t2, or t1 is a sub-type of t2
// TODO NonNull checks
//
// An object field type is a valid sub‐type if it is a Non‐Null variant of a valid sub‐type of the interface field type
// TODO: What exactly does this mean? A nullable interface field cannot be nullable for the object fields implementation??? Why not
func (builder *Builder) covariantTypeCheck(t1, t2 Type) bool {
	// An object field type is a valid sub‐type if it is equal to (the same type as) the interface field type
	if typeCheck(t1, t2) {
		return true
	}

	if t1.NonNull != t2.NonNull {
		return false
	}

	// An object field type is a valid sub‐type if it is a List type and the interface field type is also a List type and the list‐item type of the object field type is a valid sub‐type of the list‐item type of the interface field type.
	if t1.List {
		if !t2.List {
			return false
		}
		if t1.SubType != nil && t2.SubType != nil {
			return builder.covariantTypeCheck(*t1.SubType, *t2.SubType)
		}
		return false
	}

	if object, err := builder.schema.getObject(t1.Name); err == nil {
		if intrface, err := builder.schema.getInterface(t2.Name); err == nil {
			for _, objectInterface := range object.Implements {
				if objectInterface == intrface.Name {
					return true
				}
			}
		}

		if union, err := builder.schema.getUnion(t2.Name); err == nil {
			for _, unionMember := range union.Types {
				if object.Name == unionMember {
					return true
				}
			}
		}
	}
	return false
}

func findFirstDuplicate(values []string) *string {
	valueSet := make(map[string]interface{})
	for _, value := range values {
		if _, exists := valueSet[value]; exists {
			return &value
		}
		valueSet[value] = struct{}{}
	}
	return nil
}

///
// Schema Builder Helpers (mostly for readability purposes)
///

// Interfaces builds a list of Interface names
func Interfaces(interfaces ...string) []string {
	return interfaces
}

// Values builds a list of string Values
func Values(values ...string) []string {
	return values
}

// Types builds a list of Type names
func Types(types ...string) []string {
	return types
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
