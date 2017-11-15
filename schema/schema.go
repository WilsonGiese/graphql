package schema

// Schema describes the structure and behavior of a GraphQL service
type Schema struct {
	enums map[string]Enum
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
	Types       []Object
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
	enums      []Enum
	inputs     []Input
	interfaces []Interface
	objects    []Object
	scalars    []Scalar
	unions     []Union
	errors     []error
}

// NewBuilder returns a new Builder
func NewBuilder() *Builder {
	return &Builder{}
}

// Build builds and validates the Schema
func (sb *Builder) Build() (Schema, []error) {
	return Schema{}, nil
}

// Enum adds a new Enum type declaration to the Schema
func (sb *Builder) Enum(enum Enum) *Builder {
	sb.enums = append(sb.enums, enum)
	return sb
}

// Input adds a new Input type declaration to the Schema
func (sb *Builder) Input(input Input) *Builder {
	sb.inputs = append(sb.inputs, input)
	return sb
}

// Interface adds a new Interface type declaration to the Schema
func (sb *Builder) Interface(i Interface) *Builder {
	sb.interfaces = append(sb.interfaces, i)
	return sb
}

// Object adds a new Object type declaration to the Schema
func (sb *Builder) Object(object Object) *Builder {
	sb.objects = append(sb.objects, object)
	return sb
}

// Scalar adds a new Scalar type declaration to the Schema
func (sb *Builder) Scalar(scalar Scalar) *Builder {
	sb.scalars = append(sb.scalars, scalar)
	return sb
}

// Union adds a new Union type declaration to the Schema
func (sb *Builder) Union(union Union) *Builder {
	sb.unions = append(sb.unions, union)
	return sb
}

// Declare a new schema type
func (sb *Builder) Declare(declaration Declaration) *Builder {
	switch v := declaration.(type) {
	case *Enum:
		sb.Enum(*v)
	case *Input:
		sb.Input(*v)
	case *Interface:
		sb.Interface(*v)
	case *Object:
		sb.Object(*v)
	case *Scalar:
		sb.Scalar(*v)
	case *Union:
		sb.Union(*v)
	default:
		panic("Include: Invalid schema delcaration") // TODO use Declaration instead of interface{}
	}

	return sb
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
func Arguments(arguments ...Argument) (argumentsMap map[string]Argument) {
	for _, argument := range arguments {
		argumentsMap[argument.Name] = argument
	}
	return
}

// Fields builds a map from Field names to the Field itself
func Fields(fields ...Field) (fieldsMap map[string]Field) {
	for _, field := range fields {
		fieldsMap[field.Name] = field
	}
	return
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

///
// Example Schema (Test) TODO remove
///

// enum DogCommand { SIT, DOWN, HEEL }
var DogCommand = Enum{
	Name:        "DogCommand",
	Description: "Commands that a Dog may know",
	Values:      []string{"SIT", "DOWN", "HEEL"},
}

// enum CatCommand { JUMP }
var CatCommand = Enum{
	Name:        "CatCommand",
	Description: "Commands that a Cat may know",
	Values:      []string{"JUMP"},
}

// interface Sentient {
//   name: String!
// }
var Sentient = Interface{
	Name: "Sentient",
	Fields: Fields(
		Field{
			Name: "name",
			Type: NonNullStringType,
		},
	),
}

// interface Pet {
//   name: String!
// }
var Pet = Interface{
	Name: "Pet",
	Fields: Fields(
		Field{
			Name: "name",
			Type: NonNullStringType,
		},
	),
}

// type Alien implements Sentient {
//   name: String!
//   homePlanet: String
// }
var Alien = Object{
	Name:       "Alien",
	Implements: Interfaces("Sentient"),
	Fields: Fields(
		Field{
			Name:        "name",
			Description: "Name of this Alien",
			Type:        NonNullStringType,
		},
		Field{
			Name:        "homePlanet",
			Description: "The name of the planet where this Alien is from",
			Type:        StringType,
		},
	),
}

// type Human implements Sentient {
//   name: String!
// }
var Human = Object{
	Name:       "Human",
	Implements: Interfaces("Sentient"),
	Fields: Fields(
		Field{
			Name:        "name",
			Description: "Name of this Human",
			Type:        NonNullStringType,
		},
	),
}

// type Dog implements Pet {
//   name: String!
//   nickname: String
//   barkVolume: Int
//   doesKnowCommand(dogCommand: DogCommand!): Boolean!
//   isHousetrained(atOtherHomes: Boolean): Boolean!
//   owner: Human
// }
var Dog = Object{
	Name:        "Dog",
	Implements:  Interfaces("Pet"),
	Description: "Woof woof",
	Fields: Fields(
		Field{
			Name:        "name",
			Description: "Name of this Dog",
			Type:        NonNullStringType,
		},
		Field{
			Name:        "nickname",
			Description: "Nickname of this Dog",
			Type:        StringType,
		},
		Field{
			Name:        "barkVolume",
			Description: "How loud this Dog will bark",
			Type:        IntType,
		},
		Field{
			Name:        "doesKnowCommand",
			Description: "Function to determine if this Dog knows a given DogCommand",
			Type:        NonNullBooleanType,
			Arguments: Arguments(
				Argument{
					Name: "dogCommand",
					Type: DescribeNonNullType("DogCommand"),
				},
			),
		},
		Field{
			Name:        "isHousetrained",
			Description: "Function to determine if this Dog is house trained",
			Type:        NonNullBooleanType,
			Arguments: Arguments(
				Argument{
					Name: "atOtherHomes",
					Type: BooleanType,
				},
			),
		},
		Field{
			Name:        "owner",
			Description: "Owner of this dog",
			Type:        DescribeType("Human"),
		},
	),
}

// type Cat implements Pet {
//   name: String!
//   nickname: String
//   doesKnowCommand(catCommand: CatCommand!): Boolean!
//   meowVolume: Int
// }
var Cat = Object{
	Name:       "Cat",
	Implements: Interfaces("Pet"),
	Fields: Fields(
		Field{
			Name:        "name",
			Description: "Name of this Cat",
			Type:        NonNullStringType,
		},
		Field{
			Name:        "nickname",
			Description: "Nickname of this Cat",
			Type:        StringType,
		},
		Field{
			Name:        "doesKnowCommand",
			Description: "Function to determine if this Cat know a given CatCommand",
			Type:        NonNullBooleanType,
			Arguments: Arguments(
				Argument{
					Name: "catCommand",
					Type: DescribeNonNullType("CatCommand"),
				},
			),
		},
		Field{
			Name:        "meowVolume",
			Description: "How loud this cat meows",
			Type:        IntType,
		},
	),
}

// union CatOrDog = Cat | Dog
var CatOrDog = Union{
	Name:        "CatOrDog",
	Description: "A type that can either be a Cat or Dog",
	Types:       Objects(Cat, Dog),
}

// union DogOrHuman = Dog | Human
var DogOrHuman = Union{
	Name:        "DogOrHuman",
	Description: "A type that can either be a Dog or Human",
	Types:       Objects(Dog, Human),
}

// union HumanOrAlien = Human | Alien
var HumanOrAlien = Union{
	Name:        "HumanOrAlien",
	Description: "A type that can either be a Human or Alien",
	Types:       Objects(Human, Alien),
}

// type QueryRoot {
//   dog: Dog
// }
var QueryRoot = Object{
	Name:        "QueryRoot",
	Description: "The query root object for this GraphQL Schema",
	Fields: Fields(
		Field{
			Name: "dog",
			Type: DescribeType("Dog"),
		},
	),
}

// Time ...
var Time = Scalar{
	Name:        "Time",
	Description: "Represents a datetime with an ISO8601 format",
}

// The Schema's built in test are equivilent, but the first uses pre-declared
// Declarations, and the second delcares them directly in the call to Declare
func test() {
	NewBuilder().
		Scalar(Time).
		Enum(DogCommand).
		Enum(CatCommand).
		Interface(Sentient).
		Interface(Pet).
		Object(Alien).
		Object(Human).
		Object(Dog).
		Object(Cat).
		Union(CatOrDog).
		Union(DogOrHuman).
		Union(HumanOrAlien).
		Object(QueryRoot).Build()

	NewBuilder().
		Declare(Scalar{
			Name:        "Time",
			Description: "Represents a datetime with an ISO8601 format",
		}).
		Declare(Enum{
			Name:        "DogCommand",
			Description: "Commands that a Dog may know",
			Values:      []string{"SIT", "DOWN", "HEEL"},
		}).
		Declare(Enum{
			Name:        "CatCommand",
			Description: "Commands that a Cat may know",
			Values:      []string{"JUMP"},
		}).
		Declare(Interface{
			Name:        "Sentient",
			Description: "Sometimes does the thinky thinky",
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
			Name:        "Alien",
			Description: "Grey man",
			Implements:  Interfaces("Sentient"),
			Fields: Fields(
				Field{
					Name:        "name",
					Description: "Name of this Alien",
					Type:        NonNullStringType,
				},
				Field{
					Name:        "homePlanet",
					Description: "The name of the planet where this Alien is from",
					Type:        StringType,
				},
			),
		}).
		Declare(Object{
			Name:        "Human",
			Description: "Pink man",
			Implements:  Interfaces("Sentient"),
			Fields: Fields(
				Field{
					Name:        "name",
					Description: "Name of this Human",
					Type:        NonNullStringType,
				},
			),
		}).
		Declare(Object{
			Name:        "Dog",
			Description: "Woof woof",
			Implements:  Interfaces("Pet"),
			Fields: Fields(
				Field{
					Name:        "name",
					Description: "Name of this Dog",
					Type:        NonNullStringType,
				},
				Field{
					Name:        "nickname",
					Description: "Nickname of this Dog",
					Type:        StringType,
				},
				Field{
					Name:        "barkVolume",
					Description: "How loud this Dog will bark",
					Type:        IntType,
				},
				Field{
					Name:        "doesKnowCommand",
					Description: "Function to determine if this Dog knows a given DogCommand",
					Type:        NonNullBooleanType,
					Arguments: Arguments(
						Argument{
							Name: "dogCommand",
							Type: DescribeNonNullType("DogCommand"),
						},
					),
				},
				Field{
					Name:        "isHousetrained",
					Description: "Function to determine if this Dog is house trained",
					Type:        NonNullBooleanType,
					Arguments: Arguments(
						Argument{
							Name: "atOtherHomes",
							Type: BooleanType,
						},
					),
				},
				Field{
					Name:        "owner",
					Description: "Owner of this dog",
					Type:        DescribeType("Human"),
				},
			),
		}).
		Declare(Object{
			Name:        "Cat",
			Description: "Mew mew",
			Implements:  Interfaces("Pet"),
			Fields: Fields(
				Field{
					Name:        "name",
					Description: "Name of this Cat",
					Type:        NonNullStringType,
				},
				Field{
					Name:        "nickname",
					Description: "Nickname of this Cat",
					Type:        StringType,
				},
				Field{
					Name:        "doesKnowCommand",
					Description: "Function to determine if this Cat know a given CatCommand",
					Type:        NonNullBooleanType,
					Arguments: Arguments(
						Argument{
							Name: "catCommand",
							Type: DescribeNonNullType("CatCommand"),
						},
					),
				},
				Field{
					Name:        "meowVolume",
					Description: "How loud this cat meows",
					Type:        IntType,
				},
			),
		}).
		Declare(Union{
			Name:        "CatOrDog",
			Description: "A type that can either be a Cat or Dog",
			Types:       Objects(Cat, Dog),
		}).
		Declare(Union{
			Name:        "DogOrHuman",
			Description: "A type that can either be a Dog or Human",
			Types:       Objects(Dog, Human),
		}).
		Declare(Union{
			Name:        "HumanOrAlien",
			Description: "A type that can either be a Human or Alien",
			Types:       Objects(Human, Alien),
		}).
		Declare(Object{
			Name:        "QueryRoot",
			Description: "The query root object for this GraphQL Schema",
			Fields: Fields(
				Field{
					Name: "dog",
					Type: DescribeType("Dog"),
				},
			),
		}).Build()
}
