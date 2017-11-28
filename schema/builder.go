package schema

import (
	"fmt"
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

// NewBuilder returns a new Builder
func NewBuilder() *Builder {
	return &Builder{
		schema:            newSchema(),
		declaredTypeNames: make(map[string]interface{}),
	}
}

func (builder *Builder) err(format string, s ...interface{}) {
	var err error

	if len(s) == 0 {
		err = fmt.Errorf("schema validation error: %s", format)
	} else {
		err = fmt.Errorf("schema validation error: %s", fmt.Sprintf(format, s...))
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
		builder.err("Input: %s", err)
	}
	builder.schema.enums[enum.Name] = enum
	return builder
}

// Input adds a new Input type declaration to the Schema
func (builder *Builder) Input(input Input) *Builder {
	if err := builder.declareTypeName(input); err != nil {
		builder.err("Input: %s", err)
	}
	builder.schema.inputs[input.Name] = input
	return builder
}

// Interface adds a new Interface type declaration to the Schema
func (builder *Builder) Interface(intrface Interface) *Builder {
	if err := builder.declareTypeName(intrface); err != nil {
		builder.err("Interface: %s", err)
	}
	builder.schema.interfaces[intrface.Name] = intrface
	return builder
}

// Object adds a new Object type declaration to the Schema. A panic will occur
// if the Object is declared with an invalid Name, or a Type with that Name has
// already been declared
func (builder *Builder) Object(object Object) *Builder {
	if err := builder.declareTypeName(object); err != nil {
		builder.err("Object: %s", err)
	}
	builder.schema.objects[object.Name] = object
	return builder
}

// Scalar adds a new Scalar type declaration to the Schema
func (builder *Builder) Scalar(scalar Scalar) *Builder {
	if err := builder.declareTypeName(scalar); err != nil {
		builder.err("scalar: %s", err)
	}
	builder.schema.scalars[scalar.Name] = scalar
	return builder
}

// Union adds a new Union type declaration to the Schema
func (builder *Builder) Union(union Union) *Builder {
	if err := builder.declareTypeName(union); err != nil {
		builder.err("Union: %s", err)
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
		return fmt.Errorf("type with Name '%s' has already been declared", declaration.name())
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
		builder.err("Enum %s: declaration must have at least one value defined", enum.Name)
	}

	if duplicate := findFirstDuplicate(enum.Values); duplicate != nil {
		builder.err("Enum %s: declared duplicate value %s", enum.Name, *duplicate)
	}
}

// http://facebook.github.io/graphql/October2016/#sec-Input-Object-type-validation
func (builder *Builder) validateInput(input Input) {
	if len(input.Fields) == 0 {
		builder.err("Input %s: declaration must have at least one Field defined", input.Name)
	}

	for _, field := range input.Fields {
		if err := builder.validateInputField(field); err != nil {
			builder.err("Input %s(%s)", input.Name, err)
		}
	}
}

func (builder *Builder) validateInputField(field Field) error {
	if err := builder.validateName(field.Name); err != nil {
		return fmt.Errorf("Field: %s", err)
	}

	if err := builder.validateType(field.Type); err != nil {
		return fmt.Errorf("Field %s(%s)", field.Name, err)
	}

	if builder.schema.getDeclaration(field.Type).typeKind() != INPUT_OBJECT {
		return fmt.Errorf("Field %s: declared with non-Input Type %s", field.Name, field.Type)
	}

	// Input fields cannot be declared with arguments
	if len(field.Arguments) > 0 {
		return fmt.Errorf("Field %s: declared with arguments. Input fields must be declared without arguments", field.Name)
	}
	return nil
}

// http://facebook.github.io/graphql/October2016/#sec-Interface-type-validation
func (builder *Builder) validateInterface(intrface Interface) {
	if len(intrface.Fields) == 0 {
		builder.err("Interface %s: declaration must have at least one Field defined", intrface.Name)
	}

	for _, field := range intrface.Fields {
		if err := builder.validateField(field); err != nil {
			builder.err("Interface %s(%s)", intrface.Name, err)
		}
	}
}

// http://facebook.github.io/graphql/October2016/#sec-Object-type-validation
func (builder *Builder) validateObject(object Object) {

	if len(object.Fields) == 0 {
		builder.err("Object %s: declaration must have at least one Field defined", object.Name)
	}

	for _, field := range object.Fields {
		if err := builder.validateField(field); err != nil {
			builder.err("Object %s(%s)", object.Name, err)
		}
	}
}

func (builder *Builder) validateScalar(scalar Scalar) {
	// NOTE: Nothing to validate for now
}

// http://facebook.github.io/graphql/October2016/#sec-Union-type-validation
func (builder *Builder) validateUnion(union Union) {
	if len(union.Types) == 0 {
		builder.err("Union %s: declaration must have a least one member Type defined", union.Name)
	}

	for _, unionTypeName := range union.Types {
		if declaration := builder.schema.getDeclaration(DescribeType(unionTypeName)); declaration == nil {
			builder.err("Union %s: declared with unknown Type %s", union.Name, unionTypeName)
		} else if declaration.typeKind() != OBJECT {
			builder.err("Union %s: declared with non-Object Type %s", union.Name, unionTypeName)
		}
	}

	if duplicate := findFirstDuplicate(union.Types); duplicate != nil {
		builder.err("Union %s: declared duplicate Type %s", union.Name, *duplicate)
	}
}

func (builder *Builder) validateField(field Field) error {
	if err := builder.validateName(field.Name); err != nil {
		return fmt.Errorf("Field: %s", err)
	}

	if err := builder.validateType(field.Type); err != nil {
		return fmt.Errorf("Field %s(%s)", field.Name, err)
	}

	if builder.schema.getDeclaration(field.Type).typeKind() == INPUT_OBJECT {
		return fmt.Errorf("Field %s: declared with Input Type %s", field.Name, field.Type)
	}

	for _, argument := range field.Arguments {
		if err := builder.validateArgument(argument); err != nil {
			return fmt.Errorf("Field %s(%s)", field.Name, err)
		}
	}
	return nil
}

func (builder *Builder) validateArgument(argument Argument) error {
	if err := builder.validateName(argument.Name); err != nil {
		return fmt.Errorf("Argument: %s", err)
	}

	if err := builder.validateType(argument.Type); err != nil {
		return fmt.Errorf("Argument %s(%s)", argument.Name, err)
	}

	if builder.schema.getDeclaration(argument.Type).typeKind() == INPUT_OBJECT {
		return fmt.Errorf("Argument %s: declared with Input Type %s", argument.Name, argument.Type)
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
	if t.List {
		if t.SubType == nil {
			return fmt.Errorf("List Type: declared with nil SubType")
		}
		if err := builder.validateType(*t.SubType); err != nil {
			return fmt.Errorf("List Type(%s)", err)
		}
	} else {
		if err := builder.validateName(t.Name); err != nil {
			return fmt.Errorf("Type: %s", err)
		}
		if declaration := builder.schema.getDeclaration(t); declaration == nil {
			return fmt.Errorf("declared with unknown Type %s", t)
		}
	}
	return nil
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

// Interfaces builds a list of string Values
func Values(values ...string) []string {
	return values
}

// Interfaces builds a list of Type names
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
