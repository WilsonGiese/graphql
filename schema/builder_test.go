package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	Fields: Fields(
		Field{
			Name: "TestInterfaceField",
			Type: StringType,
			Arguments: Arguments(
				Argument{
					Name: "TestInterfaceFieldArgument",
					Type: BooleanType,
				},
			),
		},
	),
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
	schema := NewSchema().Build()
	assert.NotNil(t, schema)
}

///
// Valid Schema Tests
///

func TestEnumType(t *testing.T) {
	schema := NewSchema().
		Declare(TestEnum).
		Declare(Enum{
			Name:   "SingleValueEnum",
			Values: Values("JUST_ME"),
		}).
		Declare(Enum{
			Name:   "DifferentCharacterCasesEnum",
			Values: Values("lowercase", "UPPERCASE", "lowerCamelCase", "UpperCamelCase", "lower_snake_case", "UPPER_SNAKE_CASE", "withNumber1", "1", "_2", "3_4", "5_6_"),
		}).Build()
	assert.NotNil(t, schema)

	testEnum, err := schema.getEnum("TestEnum")
	assert.Nil(t, err)
	assert.Equal(t, []string{"TEST_A", "TEST_B", "TEST_C"}, testEnum.Values)

	singleValueEnum, err := schema.getEnum("SingleValueEnum")
	assert.Nil(t, err)
	assert.Equal(t, []string{"JUST_ME"}, singleValueEnum.Values)

	differentCharacterCasesEnum, err := schema.getEnum("DifferentCharacterCasesEnum")
	assert.Nil(t, err)
	assert.Equal(t, []string{"lowercase", "UPPERCASE", "lowerCamelCase", "UpperCamelCase", "lower_snake_case", "UPPER_SNAKE_CASE", "withNumber1", "1", "_2", "3_4", "5_6_"}, differentCharacterCasesEnum.Values)
}

func TestScalarType(t *testing.T) {
	schema := NewSchema().
		Declare(TestScalar).
		Declare(Scalar{
			Name: "Time_ISO8601",
		}).Build()
	assert.NotNil(t, schema)
	_, err := schema.getScalar("TestScalar")
	assert.Nil(t, err)
	_, err = schema.getScalar("Time_ISO8601")
	assert.Nil(t, err)
}

func TestInputType(t *testing.T) {
	schema := NewSchema().
		Declare(TestEnum).
		Declare(TestInput).
		Declare(Input{
			Name: "Input",
			Fields: Fields(
				Field{
					Name: "InputFieldScalar",
					Type: StringType,
				},
				Field{
					Name: "InputFieldEnum",
					Type: DescribeType("TestEnum"),
				},
				Field{
					Name: "InputFieldInput",
					Type: DescribeType("TestInput"),
				},
			),
		}).Build()
	assert.NotNil(t, schema)

	testInput, err := schema.getInput("TestInput")
	assert.Nil(t, err)
	assert.Contains(t, testInput.Fields, "TestField")

	input, err := schema.getInput("Input")
	assert.Nil(t, err)
	assert.Contains(t, input.Fields, "InputFieldScalar")
	assert.Contains(t, input.Fields, "InputFieldEnum")
	assert.Contains(t, input.Fields, "InputFieldInput")
}

func TestInterfaceType(t *testing.T) {
	schema := NewSchema().
		Declare(TestEnum).
		Declare(TestInterface).
		Declare(TestInput).
		Declare(TestObject).
		Declare(TestUnion).
		Declare(Interface{
			Name:        "Interface",
			Description: "Interface with lots of stuff",
			Fields: Fields(
				Field{
					Name: "InterfaceFieldScalar",
					Type: NonNullIDType,
				},
				Field{
					Name: "InterfaceFieldEnum",
					Type: DescribeType("TestEnum"),
				},
				Field{
					Name: "InterfaceFieldInterface",
					Type: DescribeType("TestInterface"),
				},
				Field{
					Name: "InterfaceFieldObject",
					Type: DescribeType("TestObject"),
				},
				Field{
					Name: "InterfaceFieldUnion",
					Type: DescribeType("TestUnion"),
				},
				Field{
					Name: "InterfaceFieldSelf",
					Type: DescribeType("Interface"),
				},
				Field{
					Name: "InterfaceFieldWithArguments",
					Type: IntType,
					Arguments: Arguments(
						Argument{
							Name: "Arg1",
							Type: NonNullBooleanType,
						},
						Argument{
							Name: "Arg2",
							Type: DescribeType("TestEnum"),
						},
						Argument{
							Name: "Arg3",
							Type: DescribeType("TestInput"),
						},
						Argument{
							Name:    "Arg4",
							Type:    StringType,
							Default: "DefaultString",
						},
					),
				},
			),
		}).Build()
	assert.NotNil(t, schema)

	intrface, err := schema.getInterface("Interface")
	assert.Nil(t, err)
	assert.Contains(t, intrface.Fields, "InterfaceFieldScalar")
	assert.Contains(t, intrface.Fields, "InterfaceFieldEnum")
	assert.Contains(t, intrface.Fields, "InterfaceFieldInterface")
	assert.Contains(t, intrface.Fields, "InterfaceFieldObject")
	assert.Contains(t, intrface.Fields, "InterfaceFieldUnion")
	assert.Contains(t, intrface.Fields, "InterfaceFieldSelf")
	assert.Contains(t, intrface.Fields, "InterfaceFieldWithArguments")
	assert.Contains(t, intrface.Fields["InterfaceFieldWithArguments"].Arguments, "Arg1")
	assert.Contains(t, intrface.Fields["InterfaceFieldWithArguments"].Arguments, "Arg2")
	assert.Contains(t, intrface.Fields["InterfaceFieldWithArguments"].Arguments, "Arg3")
	assert.Contains(t, intrface.Fields["InterfaceFieldWithArguments"].Arguments, "Arg4")
}

func TestSimpleObjectType(t *testing.T) {
	schema := NewSchema().
		Declare(TestEnum).
		Declare(TestInterface).
		Declare(TestInput).
		Declare(TestObject).
		Declare(TestUnion).
		Declare(Object{
			Name: "Object",
			Fields: Fields(
				Field{
					Name: "ObjectFieldEnum",
					Type: DescribeType("TestEnum"),
				},
				Field{
					Name: "ObjectFieldInterface",
					Type: DescribeType("TestInterface"),
				},
				Field{
					Name: "ObjectFieldObject",
					Type: DescribeType("TestObject"),
				},
				Field{
					Name: "ObjectFieldScalar",
					Type: BooleanType,
				},
				Field{
					Name: "ObjectFieldUnion",
					Type: DescribeType("TestUnion"),
				},
				Field{
					Name: "ObjectFieldSelf",
					Type: DescribeType("Object"),
					Arguments: Arguments(
						Argument{
							Name: "Arg1",
							Type: StringType,
						},
						Argument{
							Name: "Arg2",
							Type: DescribeType("TestEnum"),
						},
						Argument{
							Name: "Arg3",
							Type: DescribeType("TestInput"),
						},
						Argument{
							Name:    "Arg4",
							Type:    BooleanType,
							Default: true,
						},
					),
				},
			),
		}).Build()
	assert.NotNil(t, schema)

	object, err := schema.getObject("Object")
	assert.Nil(t, err)
	assert.Contains(t, object.Fields, "ObjectFieldScalar")
	assert.Contains(t, object.Fields, "ObjectFieldEnum")
	assert.Contains(t, object.Fields, "ObjectFieldInterface")
	assert.Contains(t, object.Fields, "ObjectFieldObject")
	assert.Contains(t, object.Fields, "ObjectFieldUnion")
	assert.Contains(t, object.Fields, "ObjectFieldSelf")
	assert.Contains(t, object.Fields["ObjectFieldSelf"].Arguments, "Arg1")
	assert.Contains(t, object.Fields["ObjectFieldSelf"].Arguments, "Arg2")
	assert.Contains(t, object.Fields["ObjectFieldSelf"].Arguments, "Arg3")
	assert.Contains(t, object.Fields["ObjectFieldSelf"].Arguments, "Arg4")
}

func TestObjectTypeImplementsInterfacesExactly(t *testing.T) {
	schema := NewSchema().
		Declare(TestEnum).
		Declare(TestInterface).
		Declare(TestInput).
		Declare(TestObject).
		Declare(TestUnion).
		Declare(Interface{
			Name: "Interface",
			Fields: Fields(
				Field{
					Name: "InterfaceField1",
					Type: NonNullFloatType,
				},
				Field{
					Name: "InterfaceField2",
					Type: BooleanType,
					Arguments: Arguments(
						Argument{
							Name: "InterfaceField2Arg1",
							Type: StringType,
						},
						Argument{
							Name: "InterfaceField2Arg2",
							Type: NonNullIntType,
						},
					),
				},
			),
		}).
		Declare(Object{
			Name:       "Object",
			Implements: Interfaces("Interface", "TestInterface"),
			Fields: Fields(
				Field{
					Name: "ObjectField1",
					Type: IDType,
				},
				Field{
					Name: "InterfaceField1",
					Type: NonNullFloatType,
				},
				Field{
					Name: "InterfaceField2",
					Type: BooleanType,
					Arguments: Arguments(
						Argument{
							Name: "InterfaceField2Arg1",
							Type: StringType,
						},
						Argument{
							Name: "InterfaceField2Arg2",
							Type: NonNullIntType,
						},
						Argument{
							Name: "AdditionalArg3",
							Type: FloatType,
						},
					),
				},
				Field{
					Name: "ObjectField2",
					Type: DescribeType("TestEnum"),
				},
				Field{
					Name: "TestInterfaceField",
					Type: StringType,
					Arguments: Arguments(
						Argument{
							Name: "TestInterfaceFieldArgument",
							Type: BooleanType,
						},
					),
				},
			),
		}).Build()
	assert.NotNil(t, schema)

}

func TestObjectTypeImplementsInterfaceWithCovariantTypes(t *testing.T) {
	schema := NewSchema().
		Declare(TestInterface).
		Declare(TestObject).
		Declare(TestUnion).
		Declare(Object{
			Name:       "Object1",
			Implements: Interfaces("TestInterface"),
			Fields: Fields(
				Field{
					Name: "TestInterfaceField",
					Type: StringType,
					Arguments: Arguments(
						Argument{
							Name: "TestInterfaceFieldArgument",
							Type: BooleanType,
						},
					),
				},
			),
		}).
		Declare(Interface{
			Name: "Interface",
			Fields: Fields(
				Field{
					Name: "InterfaceFieldUnion",
					Type: DescribeNonNullType("TestUnion"),
				},
				Field{
					Name: "InterfaceFieldInterface",
					Type: DescribeType("TestInterface"),
				},
				Field{
					Name: "InterfaceFieldInterfaceList",
					Type: DescribeListType(DescribeType("TestInterface")),
				},
			),
		}).
		Declare(Object{
			Name:       "Object2",
			Implements: Interfaces("Interface"),
			Fields: Fields(
				Field{
					Name: "InterfaceFieldUnion",
					Type: DescribeNonNullType("TestObject"),
				},
				Field{
					Name: "InterfaceFieldInterface",
					Type: DescribeType("Object1"),
				},
				Field{
					Name: "InterfaceFieldInterfaceList",
					Type: DescribeListType(DescribeType("Object1")),
				},
			),
		}).Build()
	assert.NotNil(t, schema)
}

func TestUnionType(t *testing.T) {
	schema := NewSchema().
		Declare(Object{
			Name: "Object1",
			Fields: Fields(
				Field{
					Name: "Field",
					Type: StringType,
				},
			),
		}).
		Declare(Object{
			Name: "Object2",
			Fields: Fields(
				Field{
					Name: "Field",
					Type: StringType,
				},
			),
		}).
		Declare(Union{
			Name:  "Object1OrObject2",
			Types: Types("Object1", "Object2"),
		}).Build()
	assert.NotNil(t, schema)
	union, err := schema.getUnion("Object1OrObject2")
	assert.Nil(t, err)
	assert.Equal(t, []string{"Object1", "Object2"}, union.Types)
}

///
// Invalid Schema Tests
///

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
	assert.Equal(t, expected, actual)
}

func TestInvalidTypeNameUndeclared(t *testing.T) {
	expected := NewValidationError("Object(TestObject) Field(TestObjectField) type declared without Name defined")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Object{
				Name: "TestObject",
				Fields: Fields(
					Field{
						Name: "TestObjectField",
						Type: DescribeType(""),
					},
				),
			}).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidTypeNameCharacters(t *testing.T) {
	expected := NewValidationError("Object(TestObject) Field(TestObjectField) type declared with an invalid Name 'Fake Type!'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Object{
				Name: "TestObject",
				Fields: Fields(
					Field{
						Name: "TestObjectField",
						Type: DescribeType("Fake Type!"),
					},
				),
			}).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidListTypeWithNilBaseType(t *testing.T) {
	expected := NewValidationError("Object(TestObject) Field(TestObjectField) type declared with a nil sub-type")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Object{
				Name: "TestObject",
				Fields: Fields(
					Field{
						Name: "TestObjectField",
						Type: Type{
							List: true,
						},
					},
				),
			}).Build()
	})
	assert.Equal(t, expected, actual)
}

///
// Invalid Enum Tests
///

func TestInvalidEnumNameUndeclared(t *testing.T) {
	expected := NewValidationError("Enum declared without Name defined")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Enum{}).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidEnumNameCharacters(t *testing.T) {
	expected := NewValidationError("Enum declared with an invalid Name 'Some Enum With Spaces'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Enum{
				Name: "Some Enum With Spaces",
			}).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidEnumWithoutValues(t *testing.T) {
	expected := NewValidationError("Enum(Test) delcared without any values defined")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Enum{
				Name: "Test",
			}).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidEnumWithDuplicateValues(t *testing.T) {
	expected := NewValidationError("Enum(Test) declared duplicate value TEST_B")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Enum{
				Name:   "Test",
				Values: Values("TEST_A", "TEST_B", "TEST_C", "TEST_B", "TEST_D"),
			}).Build()
	})
	assert.Equal(t, expected, actual)
}

///
// Invalid Input Tests
///

func TestInvalidInputNameUndeclared(t *testing.T) {
	expected := NewValidationError("Input declared without Name defined")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Input{}).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidInputNameCharacters(t *testing.T) {
	expected := NewValidationError("Input declared with an invalid Name 'InputName!'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Input{
				Name: "InputName!",
			}).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidInputWithoutFields(t *testing.T) {
	expected := NewValidationError("Input(Test) declared without any Fields defined")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Input{
				Name: "Test",
			}).Build()
	})
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
}

///
// Invalid Interface Tests
///

func TestInvalidInterfaceNameUndeclared(t *testing.T) {
	expected := NewValidationError("Interface declared without Name defined")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Interface{}).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidInterfaceNameCharacters(t *testing.T) {
	expected := NewValidationError("Interface declared with an invalid Name 'Interface%'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Interface{
				Name: "Interface%",
			}).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidInterfaceWithoutFields(t *testing.T) {
	expected := NewValidationError("Interface(Test) declared without any Fields defined")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Interface{
				Name: "Test",
			}).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidInterfaceFieldName(t *testing.T) {
	expected := NewValidationError("Interface(Test) Field declared with an invalid Name 'InvalidFieldName!'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Interface{
				Name: "Test",
				Fields: Fields(
					Field{
						Name: "InvalidFieldName!",
					},
				),
			}).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidInterfaceFieldTypeDoesNotExist(t *testing.T) {
	expected := NewValidationError("Interface(Test) Field(TestField) declared with unknown type 'FooBar'")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Interface{
				Name: "Test",
				Fields: Fields(
					Field{
						Name: "TestField",
						Type: DescribeType("FooBar"),
					},
				),
			}).Build()
	})
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
}

///
// Invalid Object Tests
///

func TestInvalidObjectNameUndeclared(t *testing.T) {
	expected := NewValidationError("Object declared without Name defined")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Object{}).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidObjectNameCharacters(t *testing.T) {
	expected := NewValidationError("Object declared with an invalid Name 'Object$'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Object{
				Name: "Object$",
			}).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidObjectWithoutFields(t *testing.T) {
	expected := NewValidationError("Object(Test) declared without any Fields defined")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Object{
				Name: "Test",
			}).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidObjectFieldName(t *testing.T) {
	expected := NewValidationError("Object(Test) Field declared with an invalid Name 'Invalid Field Name!'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Object{
				Name: "Test",
				Fields: Fields(
					Field{
						Name: "Invalid Field Name!",
					},
				),
			}).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidObjectFieldTypeDoesNotExist(t *testing.T) {
	expected := NewValidationError("Object(Test) Field(TestField) declared with unknown type 'FooBar'")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Object{
				Name: "Test",
				Fields: Fields(
					Field{
						Name: "TestField",
						Type: DescribeType("FooBar"),
					},
				),
			}).Build()
	})
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
}

func TestInvalidObjectInterfaceFieldNotImplemented(t *testing.T) {
	expected := NewValidationError("Object(Test) declared without Field(TestInterfaceField) required by Interface(TestInterface)")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Object{
				Name:       "Test",
				Implements: Interfaces("TestInterface"),
				Fields: Fields(
					Field{
						Name: "NonInterfaceField",
						Type: StringType,
					},
				),
			}).
			Declare(TestInterface).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidObjectInterfaceFieldImplementedWithWrongType(t *testing.T) {
	expected := NewValidationError("Object(Test) Field(TestInterfaceField) declared with type 'Boolean' but Interface(TestInterface) requires the type 'String' or a valid sub-type")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Object{
				Name:       "Test",
				Implements: Interfaces("TestInterface"),
				Fields: Fields(
					Field{
						Name: "NonInterfaceField",
						Type: StringType,
					},
					Field{
						Name: "TestInterfaceField",
						Type: BooleanType,
					},
				),
			}).
			Declare(TestInterface).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidObjectInterfaceFieldImplementedWithCorrectTypeButWrongVariant(t *testing.T) {
	expected := NewValidationError("Object(Test) Field(TestInterfaceField) declared with type 'String!' but Interface(TestInterface) requires the type 'String' or a valid sub-type")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Object{
				Name:       "Test",
				Implements: Interfaces("TestInterface"),
				Fields: Fields(
					Field{
						Name: "NonInterfaceField",
						Type: StringType,
					},
					Field{
						Name: "TestInterfaceField",
						Type: NonNullStringType,
					},
				),
			}).
			Declare(TestInterface).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidObjectInterfaceFieldMissingArgument(t *testing.T) {
	expected := NewValidationError("Object(Test) Field(TestInterfaceField) declared without Argument(TestInterfaceFieldArgument) required by Interface(TestInterface)")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(TestInterface).
			Declare(Object{
				Name:       "Test",
				Implements: Interfaces("TestInterface"),
				Fields: Fields(
					Field{
						Name: "NonInterfaceField",
						Type: StringType,
					},
					Field{
						Name: "TestInterfaceField",
						Type: StringType,
					},
				),
			}).
			Declare(TestInput).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidObjectInterfaceFieldArgumentImplementedWithWrongType(t *testing.T) {
	expected := NewValidationError("Object(Test) Field(TestInterfaceField) Argument(TestInterfaceFieldArgument) declared with type 'Int' but Interface(TestInterface) requires type 'Boolean'")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(TestInterface).
			Declare(Object{
				Name:       "Test",
				Implements: Interfaces("TestInterface"),
				Fields: Fields(
					Field{
						Name: "NonInterfaceField",
						Type: StringType,
					},
					Field{
						Name: "TestInterfaceField",
						Type: StringType,
						Arguments: Arguments(
							Argument{
								Name: "TestInterfaceFieldArgument",
								Type: IntType,
							},
						),
					},
				),
			}).
			Declare(TestInput).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidObjectImplementsInterfaceFieldWithAdditionalNonNullArgument(t *testing.T) {
	expected := NewValidationError("Object(Test) Field(TestInterfaceField) declared an additional Argument(AdditionalNonNullFieldArgument) with a non-null type. Since Field(TestInterfaceField) is required by Interface(TestInterface) any additional Arguments must not be required")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(TestInterface).
			Declare(Object{
				Name:       "Test",
				Implements: Interfaces("TestInterface"),
				Fields: Fields(
					Field{
						Name: "NonInterfaceField",
						Type: StringType,
					},
					Field{
						Name: "TestInterfaceField",
						Type: StringType,
						Arguments: Arguments(
							Argument{
								Name: "TestInterfaceFieldArgument",
								Type: BooleanType,
							},
							Argument{
								Name: "AdditionalFieldArgument",
								Type: IDType,
							},
							Argument{
								Name: "AdditionalNonNullFieldArgument",
								Type: NonNullIDType,
							},
						),
					},
				),
			}).
			Declare(TestInput).Build()
	})
	assert.Equal(t, expected, actual)
}

///
// Invalid Scalar Tests
///

func TestInvalidScalarNameUndeclared(t *testing.T) {
	expected := NewValidationError("Scalar declared without Name defined")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Scalar{}).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidScalarNameCharacters(t *testing.T) {
	expected := NewValidationError("Scalar declared with an invalid Name 'Scalar&'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Scalar{
				Name: "Scalar&",
			}).Build()
	})
	assert.Equal(t, expected, actual)
}

///
// Invalid Union Tests
///

func TestInvalidUnionNameUndeclared(t *testing.T) {
	expected := NewValidationError("Union declared without Name defined")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Union{}).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidUnionNameCharacters(t *testing.T) {
	expected := NewValidationError("Union declared with an invalid Name 'Un--ion'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Union{
				Name: "Un--ion",
			}).Build()
	})
	assert.Equal(t, expected, actual)
}

func TestInvalidUnionWithoutMembers(t *testing.T) {
	expected := NewValidationError("Union(Test) declared without any member types defined")

	actual := CapturePanic(func() {
		NewSchema().
			Declare(Union{
				Name: "Test",
			}).Build()
	})
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
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
	assert.Equal(t, expected, actual)
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

// GraphQL Schema Example Implementation
// http://facebook.github.io/graphql/October2016/#sec-Validation

var SampleSchema *Schema

func init() {
	SampleSchema = NewSchema().
		Declare(Enum{
			Description: "Commands that a Dog may know",
			Name:        "DogCommand",
			Values:      Values("SIT", "DOWN", "HEEL"),
		}).
		Declare(Enum{
			Description: "Commands that a Cat may know",
			Name:        "CatCommand",
			Values:      Values("JUMP"),
		}).
		Declare(Interface{
			Description: "Sometimes does the thinky thinky",
			Name:        "Sentient",
			Fields: Fields(
				Field{
					Name: "name",
					Type: NonNullStringType,
				},
			),
		}).
		Declare(Interface{
			Name: "Pet",
			Fields: Fields(
				Field{
					Name: "name",
					Type: NonNullStringType,
				},
			),
		}).
		Declare(Object{
			Description: "Sentient type Alien",
			Name:        "Alien",
			Implements:  Interfaces("Sentient"),
			Fields: Fields(
				Field{
					Description: "Name of this Alien",
					Name:        "name",
					Type:        NonNullStringType,
				},
				Field{
					Description: "The name of the planet where this Alien is from",
					Name:        "homePlanet",
					Type:        StringType,
				},
			),
		}).
		Declare(Object{
			Description: "Sentient type Human",
			Name:        "Human",
			Implements:  Interfaces("Sentient"),
			Fields: Fields(
				Field{
					Description: "Name of this Human",
					Name:        "name",
					Type:        NonNullStringType,
				},
			),
		}).
		Declare(Object{
			Description: "Pet type Dog",
			Name:        "Dog",
			Implements:  Interfaces("Pet"),
			Fields: Fields(
				Field{
					Name:        "name",
					Description: "Name of this Dog",
					Type:        NonNullStringType,
				},
				Field{
					Description: "Nickname of this Dog",
					Name:        "nickname",
					Type:        StringType,
				},
				Field{
					Description: "How loud this Dog will bark",
					Name:        "barkVolume",
					Type:        IntType,
				},
				Field{
					Description: "Function to determine if this Dog knows a given DogCommand",
					Name:        "doesKnowCommand",
					Type:        NonNullBooleanType,
					Arguments: Arguments(
						Argument{
							Name: "dogCommand",
							Type: DescribeNonNullType("DogCommand"),
						},
					),
				},
				Field{
					Description: "Function to determine if this Dog is house trained",
					Name:        "isHousetrained",
					Type:        NonNullBooleanType,
					Arguments: Arguments(
						Argument{
							Name: "atOtherHomes",
							Type: BooleanType,
						},
					),
				},
				Field{
					Description: "Owner of this dog",
					Name:        "owner",
					Type:        DescribeType("Human"),
				},
			),
		}).
		Declare(Object{
			Description: "Pet type Cat",
			Name:        "Cat",
			Implements:  Interfaces("Pet"),
			Fields: Fields(
				Field{
					Description: "Name of this Cat",
					Name:        "name",
					Type:        NonNullStringType,
				},
				Field{
					Description: "Nickname of this Cat",
					Name:        "nickname",
					Type:        StringType,
				},
				Field{
					Description: "Function to determine if this Cat know a given CatCommand",
					Name:        "doesKnowCommand",
					Type:        NonNullBooleanType,
					Arguments: Arguments(
						Argument{
							Name: "catCommand",
							Type: DescribeNonNullType("CatCommand"),
						},
					),
				},
				Field{
					Description: "How loud this cat meows",
					Name:        "meowVolume",
					Type:        IntType,
				},
			),
		}).
		Declare(Union{
			Description: "A type that can either be a Cat or Dog",
			Name:        "CatOrDog",
			Types:       Types("Cat", "Dog"),
		}).
		Declare(Union{
			Description: "A type that can either be a Dog or Human",
			Name:        "DogOrHuman",
			Types:       Types("Dog", "Human"),
		}).
		Declare(Union{
			Description: "A type that can either be a Human or Alien",
			Name:        "HumanOrAlien",
			Types:       Types("Human", "Alien"),
		}).
		Declare(Object{
			Description: "The query root object for this GraphQL Schema",
			Name:        "QueryRoot",
			Fields: Fields(
				Field{
					Name: "dog",
					Type: DescribeType("Dog"),
				},
			),
		}).Build()
}
