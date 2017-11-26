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

	// Post declaration validation order
	// Scalars
	// Enums
	// Interfaces
	// Objects
	// Unions
	// Inputs
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

func (builder *Builder) validateEnum(enum Enum) {
	if len(enum.Values) == 0 {
		builder.err("Enum %s: declaration must have at least one value defined", enum.Name)
	}
}

// http://facebook.github.io/graphql/October2016/#sec-Input-Object-type-validation
func (builder *Builder) validateInput(input Input) {
	if len(input.Fields) == 0 {
		builder.err("Input %s: declaration must have at least one Field defined", input.Name)
	}

	for _, field := range input.Fields {
		if err := builder.validateInputFieldStructure(&field); err != nil {
			builder.err("Input %s(%s)", input.Name, err)
		}
	}
}

// http://facebook.github.io/graphql/October2016/#sec-Interface-type-validation
func (builder *Builder) validateInterface(intrface Interface) {
	if len(intrface.Fields) == 0 {
		builder.err("Interface %s: declaration must have at least one Field defined", intrface.Name)
	}

	for _, field := range intrface.Fields {
		if err := builder.validateFieldStructure(&field); err != nil {
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
		if err := builder.validateFieldStructure(&field); err != nil {
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

func (builder *Builder) validateInputFieldStructure(field *Field) error {
	if err := builder.validateName(field.Name); err != nil {
		return fmt.Errorf("Field: %s", err)
	}

	if err := builder.validateTypeStructure(&field.Type); err != nil {
		return fmt.Errorf("Field %s(%s)", field.Name, err)
	}

	// Input fields cannot be declared with arguments
	if len(field.Arguments) > 0 {
		return fmt.Errorf("Field %s: declared with arguments. Input fields must be declared without arguments", field.Name)
	}
	return nil
}

func (builder *Builder) validateFieldStructure(field *Field) error {
	if err := builder.validateName(field.Name); err != nil {
		return fmt.Errorf("Field: %s", err)
	}

	if err := builder.validateTypeStructure(&field.Type); err != nil {
		return fmt.Errorf("Field %s(%s)", field.Name, err)
	}

	if len(field.Arguments) > 0 {
		for _, argument := range field.Arguments {
			if err := builder.validateArgument(&argument); err != nil {
				return fmt.Errorf("Field %s(%s)", field.Name, err)
			}
		}
	}
	return nil
}

func (builder *Builder) validateArgument(argument *Argument) error {
	if err := builder.validateName(argument.Name); err != nil {
		return fmt.Errorf("Argument: %s", err)
	}

	if err := builder.validateTypeStructure(&argument.Type); err != nil {
		return fmt.Errorf("Argument %s(%s)", argument.Name, err)
	}
	return nil
	// TODO validate default value?
}

func (builder *Builder) validateTypeStructure(t *Type) error {
	if t.List {
		if t.SubType == nil {
			return fmt.Errorf("List Type: declared with nil SubType")
		}
		if err := builder.validateTypeStructure(t.SubType); err != nil {
			return fmt.Errorf("List Type(%s)", err)
		}
	} else {
		if err := builder.validateName(t.Name); err != nil {
			return fmt.Errorf("Type: %s", err)
		}
	}
	return nil
}

func (builder *Builder) validateFieldType(field Field) error {
	declaration := builder.schema.getDeclaration(field.Type)
	if declaration == nil {
		return fmt.Errorf("Field %s: declared with invalid Type %s", field.Name, &field.Type)
	}
	if declaration.typeKind() == INPUT_OBJECT {
		return fmt.Errorf("Field %s: declared with Input Type %s", field.Name, &field.Type)
	}
	return nil
}

func (builder *Builder) validateInputFieldType(field *Field) error {
	declaration := builder.schema.getDeclaration(field.Type)
	if declaration == nil {
		return fmt.Errorf("Field %s: declared with invalid Type %s", field.Name, &field.Type)
	}
	if declaration.typeKind() == INPUT_OBJECT {
		return fmt.Errorf("Field %s: declared with non-Input Type %s", field.Name, &field.Type)
	}
	return nil
}

// Object Validation Rules:
//
// // An Object type must define one or more fields.
// // The fields of an Object type must have unique names within that Object type; no two fields may share the same name.
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
// func (builder *Builder) validateObject(object *Object, schema *Schema) {
// 		// An object type must be a super‐set of all interfaces it implements:
// 		for _, intrfaceName := range object.Implements {
// 			if intrface, exists := builder.schema.interfaces[intrfaceName]; exists {
// 				for _, intrfaceField := range intrface.Fields {
// 					if objectField, exists := object.Fields[intrfaceField.Name]; exists {
// 						// The object field must be of a type which is equal to or a sub‐type of the interface field (covariant).
//
// 					} else {
// 						builder.err("Object %s: declaration must include Field %s from Interface %s", object.Name, intrfaceField.Name, intrfaceName)
// 					}
// 				}
// 					//objectField :=
// 		}
// }

func (builder *Builder) typeCheck(t1 *Type, t2 *Type) {

}

// An object field type is a valid sub‐type if it is equal to (the same type as) the interface field type.
// An object field type is a valid sub‐type if it is an Object type and the interface field type is either an Interface type or a Union type and the object field type is a possible type of the interface field type.
// An object field type is a valid sub‐type if it is a List type and the interface field type is also a List type and the list‐item type of the object field type is a valid sub‐type of the list‐item type of the interface field type.
// An object field type is a valid sub‐type if it is a Non‐Null variant of a valid sub‐type of the interface field type.
// func (builder *Builder) covariantTypeCheckFields(t1, t2 *Type) bool {
//
// 	if !reflect.DeepEqual(t1.Type, t2.Type) {
//
// 	}
// }

// for _, interfaceName := object.Implements {
// 	if intrface, exists := schema.interfaces[interfaceName]; exists {
//
// 	} else {
// 		builder.err("Object %s declared to implement unknown Interface %s", object.Name, interfaceName)
// 	}
// }
//}

//The member types of a Union type must all be Object base types; Scalar, Interface and Union types may not be member types of a Union. Similarly, wrapping types may not be member types of a Union.
// A Union type must define one or more member types.
// func (builder *Builder) validateUnion(union *Union, schema *Schema) {
// 	if len(union.Types) == 0 {
// 		builder.err("Union %s was declared without any member types defined")
// 	}
//
// 	// All member types must be Object types
// 	for _, member := range union.Types {
// 		if _, exists := schema.objects[member]; !exists {
// 			builder.err("Union %s was declared with the member %s which is not an Object type", union.Name, member)
// 		}
// 	}
// }

///
// Schema Builder Helpers (mostly for readability purposes)
///

// Interfaces builds a list of Interfaces
func Interfaces(interfaces ...string) []string {
	return interfaces
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
