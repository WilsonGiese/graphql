package graphql

import (
	"fmt"

	schema "github.com/WilsonGiese/graphql/schema"
)

type validator struct {
	schema         *schema.Schema
	document       *Document
	errors         []error
	operationNames map[string]struct{} // Set of known Operation names
	fragmentNames  map[string]struct{} // Set of known Fragment names
}

var EXISTS struct{}

func validate(schema *schema.Schema, document *Document) []error {
	v := validator{
		schema:         schema,
		document:       document,
		operationNames: make(map[string]struct{}),
		fragmentNames:  make(map[string]struct{}),
	}

	for _, operation := range document.Operations {
		v.validateOperation(operation)
	}

	for _, fragment := range document.Fragments {
		v.validateFragment(fragment)
	}

	// Final validation checks
	//
	//  Fragments Must Be Used:
	//    Every fragment must be used at least once
	//
	return v.errors
}

// Operation Rules
//  Operation Name Uniqueness:
//    Every operation must have a unique name, even if they have differing
//    operation types (e.g. query & mutation)
//
//  Lone Anonymous Operation:
//    If any anonymous operation (namless) exists, no other operation can be
//    defined
func (v *validator) validateOperation(operation Operation) {
	// Lone Anonymous Operation
	if operation.Name == "" {
		if len(v.document.Operations) > 1 {
			v.error("Lone Anonymous Operation error: more than one operation defined with anonymous operation")
		}
	}

	// Operation Name Uniqueness
	if _, exists := v.operationNames[operation.Name]; exists {
		v.error("Operation Name Uniqueness error: duplicate operation definition found: %s", operation.Name)
	} else {
		v.operationNames[operation.Name] = EXISTS
	}

	if operation.Type == "query" {
		queryRootType := v.schema.GetDeclaration(schema.DescribeType("QueryRoot"))
		v.validateSelectionSet(queryRootType, operation.SelectionSet)
	}
}

func (v *validator) validateSelectionSet(selectionType schema.Declaration, selectionSet SelectionSet) {
	switch d := selectionType.(type) {

	case schema.Interface:
		v.validateInterfaceSelectionSet(d, selectionSet)
	case schema.Object:
		v.validateObjectSelectionSet(d, selectionSet)
	case schema.Union:
		v.validateUnionSelectionSet(d, selectionSet)
	case schema.Enum:
		if !selectionSet.IsEmpty() {
			v.error("Field Selection error: subselection not allowed on Enum '%s'", selectionType.GetName())
		}
	case schema.Scalar:
		if !selectionSet.IsEmpty() {
			v.error("Field Selection error: subselection not allowed on Scalar '%s'", selectionType.GetName())
		}
	default:
		panic("unreachable")
	}
}

func (v *validator) validateField(selectionType schema.Declaration, field Field) {
	v.validateSelectionSet(selectionType, field.SelectionSet)

	//v.validateDirectives(field.Directives)
}

func (v *validator) validateInlineFragment(inlineFragment InlineFragment) {
	if declaration := v.schema.GetDeclaration(schema.DescribeType(inlineFragment.Type)); declaration == nil {
		v.error("Inline Fragment Spread Type Existence error: target type '%s' does not exist in the schema", inlineFragment.Type)
	} else {
		v.validateSelectionSet(declaration, inlineFragment.SelectionSet)
	}

	// TODO nested fragment validation
}

func (v *validator) validateFragmentSpread(fragmentSpread FragmentSpread) {
	if _, err := v.document.GetFragment(fragmentSpread.Name); err != nil {
		v.error("Fragment Spread error: Fragment '%s' is not defined", fragmentSpread.Name)
	}
	// TODO valid parent type spread & nesred ty
}

func (v *validator) validateInterfaceSelectionSet(intrface schema.Interface, selectionSet SelectionSet) {
	if selectionSet.IsEmpty() {
		v.error("Field Selection error: Interface type '%s' must have a subselection", intrface.Name)
	} else {
		for _, selectedField := range selectionSet.Fields {
			if actualField, exists := intrface.Fields[selectedField.Name]; exists {
				v.validateField(v.schema.GetDeclaration(actualField.Type), selectedField)
				v.validateArguments(actualField, selectedField)
			} else {
				v.error("Field Selection error: Interface type '%s' does not contain the field '%s'", intrface.Name, selectedField.Name)
			}
		}

		for _, inlineFragment := range selectionSet.InlineFragments {
			v.validateInlineFragment(inlineFragment)
		}

		for _, fragmentSpread := range selectionSet.FragmentSpreads {
			v.validateFragmentSpread(fragmentSpread)
		}
	}
}

func (v *validator) validateObjectSelectionSet(object schema.Object, selectionSet SelectionSet) {
	if selectionSet.IsEmpty() {
		v.error("Field Selection error: Object type '%s' must have a subselection", object.Name)
	} else {
		for _, selectedField := range selectionSet.Fields {
			if actualField, exists := object.Fields[selectedField.Name]; exists {
				v.validateField(v.schema.GetDeclaration(actualField.Type), selectedField)
				v.validateArguments(actualField, selectedField)
			} else {
				v.error("Field Selection error: Object type '%s' does not contain the field '%s'", object.Name, selectedField.Name)
			}
		}

		for _, inlineFragment := range selectionSet.InlineFragments {
			v.validateInlineFragment(inlineFragment)
		}

		for _, fragmentSpread := range selectionSet.FragmentSpreads {
			v.validateFragmentSpread(fragmentSpread)
		}
	}
}

func (v *validator) validateUnionSelectionSet(union schema.Union, selectionSet SelectionSet) {
	if selectionSet.IsEmpty() {
		v.error("Field Selection error: Union type '%s' must have a subselection", union.Name)
	} else {
		for _, selectedField := range selectionSet.Fields {
			if selectedField.Name != "__typename" {
				v.error("Field Selection error: cannot select non-metadata field from Union '%s'. Use fragment spreads to select fields from Union member types", union.Name)
			}
		}

		for _, inlineFragment := range selectionSet.InlineFragments {
			v.validateInlineFragment(inlineFragment)
		}

		for _, fragmentSpread := range selectionSet.FragmentSpreads {
			v.validateFragmentSpread(fragmentSpread)
		}
	}
}

func (v *validator) validateArguments(field schema.Field, selectedField Field) {
	for argName := range selectedField.Arguments {
		if _, exists := field.Arguments[argName]; !exists {
			v.error("Field Argument error: provided invalid argument '%s' to field '%s'", argName, field.Name)
		}
	}
}

// Fragment Rules
//  Fragment Name Uniqueness
//    Every fragment must have a unique name
//  Fragment Spread Type Existence TODO must validated for inline fragments as well
//    The "on" type for the Fragment definition must exist in the Schema
//
func (v *validator) validateFragment(fragment Fragment) {
	// Fragment Name Uniqueness
	if _, exists := v.fragmentNames[fragment.Name]; exists {
		v.error("Fragment Name Uniqueness error: duplicate fragment definition found '%s'", fragment.Name)
	} else {
		v.fragmentNames[fragment.Name] = EXISTS
	}

	// Fragment Spread Type Existence
	if declaration := v.schema.GetDeclaration(schema.DescribeType(fragment.Type)); declaration == nil {
		v.error("Fragment Spread Type Existence error: target type '%s' does not exist in the schema", fragment.Type)
	} else {
		v.validateSelectionSet(declaration, fragment.SelectionSet)
	}

	//v.validateDirectives(fragment.Directives)
}

func (v *validator) error(format string, s ...interface{}) {
	v.errors = append(v.errors, fmt.Errorf(format, s...))
}
