package schema

import (
	"testing"
)

var TestEnum = Enum{
	Name:   "TestEnum",
	Values: Values("TEST_A", "TEST_B", "TEST_C"),
}

var TestInput = Input{
	Name: "TestInput",
	Fields: Fields(Field{
		Name: "TestField",
		Type: IntType,
	}),
}

var TestInterface = Interface{
	Name: "TestInterface",
	Fields: Fields(Field{
		Name: "TestInterfaceField",
		Type: BooleanType,
	}),
}

var TestObject = Object{
	Name: "TestObject",
	Fields: Fields(Field{
		Name: "TestField",
		Type: StringType,
	}),
}

var TestScalar = Scalar{
	Name: "TestScalar",
}

var TestUnion = Union{
	Name:  "TestUnion",
	Types: Types("TestObject"),
}

// Tests if default Schema can be built successfully, should panic otherwise
func TestBuildDefaultSchema(t *testing.T) {
	NewSchema().Build()
}

func TestInvalidSchemaMultipleDeclarationsWithTheSameName(t *testing.T) {
	expected := NewValidationError("Object declared with name 'Duplicate' but another type with that name has already been declared")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Enum{
				Name: "Test",
			}).
			Declare(Union{
				Name: "Duplicate",
			}).
			Declare(Union{
				Name: "Type",
			}).
			Declare(Object{
				Name: "Duplicate",
			}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

///
// Invalid Enum Tests
///

func TestInvalidEnumNameUndeclared(t *testing.T) {
	expected := NewValidationError("Enum declared without Name defined")

	actual := CapturePanic(func() {
		NewSchema().Declare(Enum{}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidEnumNameCharacters(t *testing.T) {
	expected := NewValidationError("Enum declared with an invalid Name 'Some Enum With Spaces'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().Declare(Enum{
			Name: "Some Enum With Spaces",
		}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidEnumWithoutValues(t *testing.T) {
	expected := NewValidationError("Enum(Test) delcared without any values defined")

	actual := CapturePanic(func() {
		NewSchema().Declare(Enum{
			Name: "Test",
		}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidEnumWithDuplicateValues(t *testing.T) {
	expected := NewValidationError("Enum(Test) declared duplicate value TEST_B")

	actual := CapturePanic(func() {
		NewSchema().Declare(Enum{
			Name:   "Test",
			Values: Values("TEST_A", "TEST_B", "TEST_C", "TEST_B", "TEST_D"),
		}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

///
// Invalid Input Tests
///

func TestInvalidInputNameUndeclared(t *testing.T) {
	expected := NewValidationError("Input declared without Name defined")

	actual := CapturePanic(func() {
		NewSchema().Declare(Input{}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidInputNameCharacters(t *testing.T) {
	expected := NewValidationError("Input declared with an invalid Name 'InputName!'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().Declare(Input{
			Name: "InputName!",
		}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidInputWithoutFields(t *testing.T) {
	expected := NewValidationError("Input(Test) declared without any Fields defined")

	actual := CapturePanic(func() {
		NewSchema().Declare(Input{
			Name: "Test",
		}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidInputFieldName(t *testing.T) {
	expected := NewValidationError("Input(Test) Field declared with an invalid Name 'InvalidFieldName!'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Input{
				Name: "Test",
				Fields: Fields(Field{
					Name: "InvalidFieldName!",
				}),
			}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidInputFieldTypeDoesNotExist(t *testing.T) {
	expected := NewValidationError("Input(Test) Field(TestField) declared with unknown type 'FooBar'")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Input{
				Name: "Test",
				Fields: Fields(Field{
					Name: "TestField",
					Type: DescribeType("FooBar"),
				}),
			}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidInputFieldUnacceptableInterfaceType(t *testing.T) {
	expected := NewValidationError("Input(Test) Field(TestInterfaceField) declared with invalid Type 'TestInterface'. An Input Field type must be Input, Scalar, or Enum")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(TestInterface).
			Declare(Input{
				Name: "Test",
				Fields: Fields(
					Field{
						Name: "TestInterfaceField",
						Type: DescribeType("TestInterface"),
					},
				),
			}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidInputFieldUnacceptableObjectType(t *testing.T) {
	expected := NewValidationError("Input(Test) Field(TestObjectField) declared with invalid Type 'TestObject'. An Input Field type must be Input, Scalar, or Enum")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(TestObject).
			Declare(Input{
				Name: "Test",
				Fields: Fields(
					Field{
						Name: "TestObjectField",
						Type: DescribeType("TestObject"),
					},
				),
			}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidInputFieldUnacceptableUnionType(t *testing.T) {
	expected := NewValidationError("Input(Test) Field(TestUnionField) declared with invalid Type 'TestUnion'. An Input Field type must be Input, Scalar, or Enum")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(TestObject).
			Declare(TestUnion).
			Declare(Input{
				Name: "Test",
				Fields: Fields(
					Field{
						Name: "TestUnionField",
						Type: DescribeType("TestUnion"),
					},
				),
			}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidInputFieldWithArguments(t *testing.T) {
	expected := NewValidationError("Input(Test) Field(TestField) declared with arguments. Input fields must be declared without arguments")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Input{
				Name: "Test",
				Fields: Fields(Field{
					Name: "TestField",
					Type: BooleanType,
					Arguments: Arguments(
						Argument{
							Name: "TestArgument",
							Type: FloatType,
						},
					),
				}),
			}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

///
// Invalid Interface Tests
///

func TestInvalidInterfaceNameUndeclared(t *testing.T) {
	expected := NewValidationError("Interface declared without Name defined")

	actual := CapturePanic(func() {
		NewSchema().Declare(Interface{}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidInterfaceNameCharacters(t *testing.T) {
	expected := NewValidationError("Interface declared with an invalid Name 'Interface%'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().Declare(Interface{
			Name: "Interface%",
		}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidInterfaceWithoutFields(t *testing.T) {
	expected := NewValidationError("Interface(Test) declared without any Fields defined")

	actual := CapturePanic(func() {
		NewSchema().Declare(Interface{
			Name: "Test",
		}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidInterfaceFieldName(t *testing.T) {
	expected := NewValidationError("Interface(Test) Field declared with an invalid Name 'InvalidFieldName!'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().Declare(Interface{
			Name: "Test",
			Fields: Fields(
				Field{
					Name: "InvalidFieldName!",
				},
			),
		}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidInterfaceFieldTypeDoesNotExist(t *testing.T) {
	expected := NewValidationError("Interface(Test) Field(TestField) declared with unknown type 'FooBar'")

	actual := CapturePanic(func() {
		NewSchema().Declare(Interface{
			Name: "Test",
			Fields: Fields(
				Field{
					Name: "TestField",
					Type: DescribeType("FooBar"),
				},
			),
		}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidInterfaceFieldTypeUnacceptable(t *testing.T) {
	expected := NewValidationError("Interface(Test) Field(TestField) declared with Input type 'TestInput'")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Interface{
				Name: "Test",
				Fields: Fields(
					Field{
						Name: "TestField",
						Type: DescribeType("TestInput"),
					},
				),
			}).
			Declare(TestInput).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidInterfaceFieldArgumentName(t *testing.T) {
	expected := NewValidationError("Interface(Test) Field(TestField) Argument: declared with an invalid Name 'Test Argument'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Interface{
				Name: "Test",
				Fields: Fields(
					Field{
						Name: "TestField",
						Type: StringType,
						Arguments: Arguments(
							Argument{
								Name: "Test Argument",
							},
						),
					},
				),
			}).
			Declare(TestInput).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidInterfaceFieldArgumentTypeDoesNotExist(t *testing.T) {
	expected := NewValidationError("Interface(Test) Field(TestField) Argument(TestArgument) declared with unknown type 'FooBar'")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Interface{
				Name: "Test",
				Fields: Fields(
					Field{
						Name: "TestField",
						Type: StringType,
						Arguments: Arguments(
							Argument{
								Name: "TestArgument",
								Type: DescribeType("FooBar"),
							},
						),
					},
				),
			}).
			Declare(TestInput).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidInterfaceFieldArgumentUnacceptableInterfaceType(t *testing.T) {
	expected := NewValidationError("Interface(Test) Field(TestField) Argument(TestArgumentInterface) declared with invalid type 'TestInterface'. An Argument Type must be Input, Scalar, or Enum")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(TestInterface).
			Declare(Interface{
				Name: "Test",
				Fields: Fields(
					Field{
						Name: "TestField",
						Type: StringType,
						Arguments: Arguments(
							Argument{
								Name: "TestArgumentInterface",
								Type: DescribeType("TestInterface"),
							},
						),
					},
				),
			}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidInterfaceFieldArgumentUnacceptableObjectType(t *testing.T) {
	expected := NewValidationError("Interface(Test) Field(TestField) Argument(TestArgumentObject) declared with invalid type 'TestObject'. An Argument Type must be Input, Scalar, or Enum")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(TestObject).
			Declare(Interface{
				Name: "Test",
				Fields: Fields(
					Field{
						Name: "TestField",
						Type: StringType,
						Arguments: Arguments(
							Argument{
								Name: "TestArgumentObject",
								Type: DescribeType("TestObject"),
							},
						),
					},
				),
			}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidInterfaceFieldArgumentUnacceptableUnionType(t *testing.T) {
	expected := NewValidationError("Interface(Test) Field(TestField) Argument(TestArgumentUnion) declared with invalid type 'TestUnion'. An Argument Type must be Input, Scalar, or Enum")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(TestUnion).
			Declare(Interface{
				Name: "Test",
				Fields: Fields(
					Field{
						Name: "TestField",
						Type: StringType,
						Arguments: Arguments(
							Argument{
								Name: "TestArgumentUnion",
								Type: DescribeType("TestUnion"),
							},
						),
					},
				),
			}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidInterfaceFieldArgumentNonNullWithDefaultValue(t *testing.T) {
	expected := NewValidationError("Interface(Test) Field(TestField) Argument(TestArgument) declared with a default value, but its type is non-null")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Interface{
				Name: "Test",
				Fields: Fields(
					Field{
						Name: "TestField",
						Type: StringType,
						Arguments: Arguments(
							Argument{
								Name:    "TestArgument",
								Type:    NonNullStringType,
								Default: "DefaultString",
							},
						),
					},
				),
			}).
			Declare(TestInput).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

///
// Invalid Object Tests
///

func TestInvalidObjectNameUndeclared(t *testing.T) {
	expected := NewValidationError("Object declared without Name defined")

	actual := CapturePanic(func() {
		NewSchema().Declare(Object{}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidObjectNameCharacters(t *testing.T) {
	expected := NewValidationError("Object declared with an invalid Name 'Object$'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().Declare(Object{
			Name: "Object$",
		}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidObjectWithoutFields(t *testing.T) {
	expected := NewValidationError("Object(Test) declared without any Fields defined")

	actual := CapturePanic(func() {
		NewSchema().Declare(Object{
			Name: "Test",
		}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidObjectFieldName(t *testing.T) {
	expected := NewValidationError("Object(Test) Field declared with an invalid Name 'Invalid Field Name!'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().Declare(Object{
			Name: "Test",
			Fields: Fields(
				Field{
					Name: "Invalid Field Name!",
				},
			),
		}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidObjectFieldTypeDoesNotExist(t *testing.T) {
	expected := NewValidationError("Object(Test) Field(TestField) declared with unknown type 'FooBar'")

	actual := CapturePanic(func() {
		NewSchema().Declare(Object{
			Name: "Test",
			Fields: Fields(
				Field{
					Name: "TestField",
					Type: DescribeType("FooBar"),
				},
			),
		}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidObjectFieldTypeUnacceptable(t *testing.T) {
	expected := NewValidationError("Object(Test) Field(TestField) declared with Input type 'TestInput'")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Object{
				Name: "Test",
				Fields: Fields(
					Field{
						Name: "TestField",
						Type: DescribeType("TestInput"),
					},
				),
			}).
			Declare(TestInput).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidObjectFieldArgumentName(t *testing.T) {
	expected := NewValidationError("Object(Test) Field(TestField) Argument: declared with an invalid Name 'Test Argument'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Object{
				Name: "Test",
				Fields: Fields(
					Field{
						Name: "TestField",
						Type: StringType,
						Arguments: Arguments(
							Argument{
								Name: "Test Argument",
							},
						),
					},
				),
			}).
			Declare(TestInput).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidObjectFieldArgumentTypeDoesNotExist(t *testing.T) {
	expected := NewValidationError("Object(Test) Field(TestField) Argument(TestArgument) declared with unknown type 'FooBar'")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Object{
				Name: "Test",
				Fields: Fields(
					Field{
						Name: "TestField",
						Type: StringType,
						Arguments: Arguments(
							Argument{
								Name: "TestArgument",
								Type: DescribeType("FooBar"),
							},
						),
					},
				),
			}).
			Declare(TestInput).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidObjectFieldArgumentUnacceptableInterfaceType(t *testing.T) {
	expected := NewValidationError("Object(Test) Field(TestField) Argument(TestArgumentInterface) declared with invalid type 'TestInterface'. An Argument Type must be Input, Scalar, or Enum")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(TestInterface).
			Declare(Object{
				Name: "Test",
				Fields: Fields(
					Field{
						Name: "TestField",
						Type: StringType,
						Arguments: Arguments(
							Argument{
								Name: "TestArgumentInterface",
								Type: DescribeType("TestInterface"),
							},
						),
					},
				),
			}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidObjectFieldArgumentUnacceptableObjectType(t *testing.T) {
	expected := NewValidationError("Object(Test) Field(TestField) Argument(TestArgumentObject) declared with invalid type 'TestObject'. An Argument Type must be Input, Scalar, or Enum")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(TestObject).
			Declare(Object{
				Name: "Test",
				Fields: Fields(
					Field{
						Name: "TestField",
						Type: StringType,
						Arguments: Arguments(
							Argument{
								Name: "TestArgumentObject",
								Type: DescribeType("TestObject"),
							},
						),
					},
				),
			}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidObjectFieldArgumentUnacceptableUnionType(t *testing.T) {
	expected := NewValidationError("Object(Test) Field(TestField) Argument(TestArgumentUnion) declared with invalid type 'TestUnion'. An Argument Type must be Input, Scalar, or Enum")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(TestUnion).
			Declare(Object{
				Name: "Test",
				Fields: Fields(
					Field{
						Name: "TestField",
						Type: StringType,
						Arguments: Arguments(
							Argument{
								Name: "TestArgumentUnion",
								Type: DescribeType("TestUnion"),
							},
						),
					},
				),
			}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidObjectFieldArgumentNonNullWithDefaultValue(t *testing.T) {
	expected := NewValidationError("Object(Test) Field(TestField) Argument(TestArgument) declared with a default value, but its type is non-null")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Object{
				Name: "Test",
				Fields: Fields(
					Field{
						Name: "TestField",
						Type: StringType,
						Arguments: Arguments(
							Argument{
								Name:    "TestArgument",
								Type:    NonNullStringType,
								Default: "DefaultString",
							},
						),
					},
				),
			}).
			Declare(TestInput).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidObjectImplementsUnknownInterface(t *testing.T) {
	expected := NewValidationError("Object(Test) declared implementing unknown Interface 'FooBar'")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Object{
				Name:       "Test",
				Implements: Interfaces("FooBar"),
				Fields: Fields(
					Field{
						Name: "NonInterfaceField",
						Type: StringType,
					},
				),
			}).
			Declare(TestInput).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidObjectInterfaceFieldNotImplemented(t *testing.T) {
	expected := NewValidationError("Object(Test) declared without Field(SomeInterfaceField1) required by Interface(SomeInterface)")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Interface{
				Name: "SomeInterface",
				Fields: Fields(
					Field{
						Name: "SomeInterfaceField1",
						Type: StringType,
					},
				),
			}).
			Declare(Object{
				Name:       "Test",
				Implements: Interfaces("SomeInterface"),
				Fields: Fields(
					Field{
						Name: "NonInterfaceField",
						Type: StringType,
					},
				),
			}).
			Declare(TestInput).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidObjectInterfaceFieldImplementedWithWrongType(t *testing.T) {
	expected := NewValidationError("Object(Test) Field(SomeInterfaceField1) declared with type 'Boolean' but Interface(SomeInterface) requires the type 'String' or a valid sub-type")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Interface{
				Name: "SomeInterface",
				Fields: Fields(
					Field{
						Name: "SomeInterfaceField1",
						Type: StringType,
					},
				),
			}).
			Declare(Object{
				Name:       "Test",
				Implements: Interfaces("SomeInterface"),
				Fields: Fields(
					Field{
						Name: "NonInterfaceField",
						Type: StringType,
					},
					Field{
						Name: "SomeInterfaceField1",
						Type: BooleanType,
					},
				),
			}).
			Declare(TestInput).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidObjectInterfaceFieldImplementedWithCorrectTypeButWrongVariant(t *testing.T) {
	expected := NewValidationError("Object(Test) Field(SomeInterfaceField1) declared with type 'String!' but Interface(SomeInterface) requires the type 'String' or a valid sub-type")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Interface{
				Name: "SomeInterface",
				Fields: Fields(
					Field{
						Name: "SomeInterfaceField1",
						Type: StringType,
					},
				),
			}).
			Declare(Object{
				Name:       "Test",
				Implements: Interfaces("SomeInterface"),
				Fields: Fields(
					Field{
						Name: "NonInterfaceField",
						Type: StringType,
					},
					Field{
						Name: "SomeInterfaceField1",
						Type: NonNullStringType,
					},
				),
			}).
			Declare(TestInput).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

///
// Invalid Scalar Tests
///

func TestInvalidScalarNameUndeclared(t *testing.T) {
	expected := NewValidationError("Scalar declared without Name defined")

	actual := CapturePanic(func() {
		NewSchema().Declare(Scalar{}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidScalarNameCharacters(t *testing.T) {
	expected := NewValidationError("Scalar declared with an invalid Name 'Scalar&'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().Declare(Scalar{
			Name: "Scalar&",
		}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

///
// Invalid Union Tests
///

func TestInvalidUnionNameUndeclared(t *testing.T) {
	expected := NewValidationError("Union declared without Name defined")

	actual := CapturePanic(func() {
		NewSchema().Declare(Union{}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidUnionNameCharacters(t *testing.T) {
	expected := NewValidationError("Union declared with an invalid Name 'Union()'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().Declare(Union{
			Name: "Union()",
		}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidUnionWithoutMembers(t *testing.T) {
	expected := NewValidationError("Union(Test) declared without any member types defined")

	actual := CapturePanic(func() {
		NewSchema().Declare(Union{
			Name: "Test",
		}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidUnionWithDuplicateMembers(t *testing.T) {
	expected := NewValidationError("Union(Test) declared duplicate type TestObject")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Union{
				Name:  "Test",
				Types: Types("TestObject", "TestObject"),
			}).
			Declare(TestObject).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidUnionMemberTypeDoesNotExist(t *testing.T) {
	expected := NewValidationError("Union(TestUnion) declared with unknown type TestUnknownObject")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Union{
				Name:  "TestUnion",
				Types: Types("TestObject", "TestUnknownObject"),
			}).
			Declare(TestObject).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidUnionMemberNonObjectMemberType(t *testing.T) {
	expected := NewValidationError("Union(TestUnion) declared with member type TestEnum. Union members must be Objects")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Union{
				Name:  "TestUnion",
				Types: Types("TestEnum", "TestInput", "TestInterface", "TestScalar", "TestObject"),
			}).
			Declare(TestEnum).
			Declare(TestInput).
			Declare(TestInterface).
			Declare(TestObject).
			Declare(TestScalar).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

// Capture the panic from function f and return it as an error
func CapturePanic(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		} else {
			err = nil
		}
	}()
	f()
	return
}

func AssertSchemaValidationError(expected error, actual error, t *testing.T) {
	if actual == nil {
		t.Errorf("expected '%s' but error was nil", expected)
	} else {
		if actual.Error() != expected.Error() {
			t.Errorf("expected '%s' but got '%s'", expected, actual)
		}
	}
}
