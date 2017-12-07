package graphql

import "fmt"

type validator struct {
	schema         *Schema
	document       *Document
	errors         []error
	operationNames map[string]struct{} // Set of known Operation names
	fragmentNames  map[string]struct{} // Set of known Fragment names
}

var EXISTS struct{}

func validate(schema *Schema, document *Document) {
	v := validator{
		schema:   schema,
		document: document,
	}

	for i := 0; i < len(document.Operations); i++ {
		v.validateOperation(&document.Operations[i])
	}

	for i := 0; i < len(document.Fragments); i++ {
		v.validateFragment(&document.Fragments[i])
	}

  // Final validation checks
  //
  //  Fragments Must Be Used:
  //    Every fragment must be used at least once
  //
}

// Operation Rules
//  Operation Name Uniqueness:
//    Every operation must have a unique name, even if they have differing
//    operation types (e.g. query & mutation)
//
//  Lone Anonymous Operation:
//    If any anonymous operation (namless) exists, no other operation can be
//    defined
func (v *validator) validateOperation(operation *Operation) {
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
}

// Fragment Rules
//  Fragment Name Uniqueness
//    Every fragment must have a unique name
//  Fragment Spread Type Existence TODO must validated for inline fragments as well
//    The "on" type for the Fragment definition must exist in the Schema
//
func (v *validator) validateFragment(fragment *Fragment) {
  // Fragment Name Uniqueness
  if _, exists := v.operationNames[operation.Name]; exists {
    v.error("Fragment Name Uniqueness error: duplicate fragment definition found: %s", operation.Name)
  } else {
    v.fragmentNames[fragment.Name] = EXISTS
  }

  // Fragment Spread Type Existence
  if v.schema.
}

func (v *validator) error(format string, s ...interface{}) {
	v.errors = append(v.errors, fmt.Errorf(format, s))
}
