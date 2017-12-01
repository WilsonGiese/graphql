package schema

import (
	"testing"
)

// Tests if default Schema can be built successfully, should panic otherwise
func TestBuildDefaultSchema(t *testing.T) {
	NewSchema().Build()
}

func TestInvalidSchemaMultipleDeclarationsWithTheSameName(t *testing.T) {
	expected := NewSchemaValidationError("Object declared with Name 'Duplicate' but another type with that name has already been declared")

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
	expected := NewSchemaValidationError("Enum declared without Name defined")

	actual := CapturePanic(func() {
		NewSchema().Declare(Enum{}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidEnumNameCharacters(t *testing.T) {
	expected := NewSchemaValidationError("Enum declared with an invalid Name 'Some Enum With Spaces'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().Declare(Enum{
			Name: "Some Enum With Spaces",
		}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidEnumWithoutValues(t *testing.T) {
	expected := NewSchemaValidationError("Enum(Test) delcared without any values defined")

	actual := CapturePanic(func() {
		NewSchema().Declare(Enum{
			Name: "Test",
		}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidEnumWithDuplicateValues(t *testing.T) {
	expected := NewSchemaValidationError("Enum(Test) declared duplicate value TEST_B")

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
	expected := NewSchemaValidationError("Input declared without Name defined")

	actual := CapturePanic(func() {
		NewSchema().Declare(Input{}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidInputNameCharacters(t *testing.T) {
	expected := NewSchemaValidationError("Input declared with an invalid Name 'InputName!'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().Declare(Input{
			Name: "InputName!",
		}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidInputWithoutFields(t *testing.T) {
	expected := NewSchemaValidationError("Input(Test) declared without any Fields defined")

	actual := CapturePanic(func() {
		NewSchema().Declare(Input{
			Name: "Test",
		}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

///
// Invalid Interface Tests
///

func TestInvalidInterfaceNameUndeclared(t *testing.T) {
	expected := NewSchemaValidationError("Interface declared without Name defined")

	actual := CapturePanic(func() {
		NewSchema().Declare(Interface{}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidInterfaceNameCharacters(t *testing.T) {
	expected := NewSchemaValidationError("Interface declared with an invalid Name 'Interface%'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().Declare(Interface{
			Name: "Interface%",
		}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

///
// Invalid Object Tests
///

func TestInvalidObjectNameUndeclared(t *testing.T) {
	expected := NewSchemaValidationError("Object declared without Name defined")

	actual := CapturePanic(func() {
		NewSchema().Declare(Object{}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidObjectNameCharacters(t *testing.T) {
	expected := NewSchemaValidationError("Object declared with an invalid Name 'Object$'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().Declare(Object{
			Name: "Object$",
		}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

///
// Invalid Scalar Tests
///

func TestInvalidScalarNameUndeclared(t *testing.T) {
	expected := NewSchemaValidationError("Scalar declared without Name defined")

	actual := CapturePanic(func() {
		NewSchema().Declare(Scalar{}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidScalarNameCharacters(t *testing.T) {
	expected := NewSchemaValidationError("Scalar declared with an invalid Name 'Scalar&'. A Name must only consist of ASCII letters, numbers, and underscores")

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
	expected := NewSchemaValidationError("Union declared without Name defined")

	actual := CapturePanic(func() {
		NewSchema().Declare(Union{}).Build()
	})
	AssertSchemaValidationError(expected, actual, t)
}

func TestInvalidUnionNameCharacters(t *testing.T) {
	expected := NewSchemaValidationError("Union declared with an invalid Name 'Union()'. A Name must only consist of ASCII letters, numbers, and underscores")

	actual := CapturePanic(func() {
		NewSchema().Declare(Union{
			Name: "Union()",
		}).Build()
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
		t.Errorf("Expected error '%s' but error was nil", expected)
	} else {
		if actual.Error() != expected.Error() {
			t.Errorf("Expected error '%s' but got '%s'", expected, actual)
		}
	}
}
