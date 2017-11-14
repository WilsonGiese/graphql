package graphql

// Schema describes the structure and behavior of a GraphQL service
type Schema struct {
}

// Scalar describes a scalar type within a Schema
type Scalar struct {
	Name        string
	Description string
}

// EnumSchema describes an enum type within a Schema
type EnumSchema struct {
	Name        string
	Values      []string
	Description string
}

// InputSchema describes an input type object within a Schema
type InputSchema struct {
	Name        string
	Description string
}

// Fields represents a map from field names (strings) to a TypeSchema
type Fields map[string]TypeSchema

// InterfaceSchema descrives an object type interface within a Schema
type InterfaceSchema struct {
	Name        string
	Description string
	Fields      Fields
}

// UnionSchema describes a union of Types within a Schema
type UnionSchema struct {
	Name        string
	Description string
	Types       []ObjectSchema
}

// ObjectSchema describes an Object Type defined within a Schema
type ObjectSchema struct {
	Name        string
	Implements  InterfaceSchema
	Description string
	Fields      ObjectFields
}

// ObjectFieldSchema describes a Field for a Type defined within a Schema
type ObjectFieldSchema struct {
	Name        string
	Description string
	Type        TypeSchema
	Arguments   Arguments
}

// ObjectFields represents a map of field names (strings) to an ObjectFieldSchema
type ObjectFields map[string]ObjectFieldSchema

// TypeSchema represents a Type in a Schema
type TypeSchema struct {
	Name    string
	NonNull bool
	List    bool
	SubType *TypeSchema
}

// ArgumentSchema defines an argument for a Field defined within a Type
type ArgumentSchema struct {
	Name    string
	Type    TypeSchema
	Default interface{}
}

// Arguments represents a map from argument names (strings) to an ArgumentSchema
type Arguments map[string]ArgumentSchema

// DescribeType returns a TypeSchema for the given type which can be null
func DescribeType(name string) (t TypeSchema) {
	t.Name = name
	return
}

// DescribeNonNullType returns a non-null TypeSchema for the given type
func DescribeNonNullType(name string) (t TypeSchema) {
	t.Name = name
	t.NonNull = true
	return
}

// DescribeListType returns a TypeSchema for the given subType which can be null
func DescribeListType(subType TypeSchema) (t TypeSchema) {
	t.List = true
	t.SubType = &subType
	return
}

// DescribeNonNullListType returns a non-null TypeSchema for the given subType
func DescribeNonNullListType(subType TypeSchema) (t TypeSchema) {
	t.List = true
	t.NonNull = true
	t.SubType = &subType
	return
}

// Common Pre-defined Types (Int, Float, String, Boolean, ID)

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
// Example Schema (Test)
///

// enum DogCommand { SIT, DOWN, HEEL }
var DogCommand = EnumSchema{
	Name:        "DogCommand",
	Description: "Commands that a Dog may know",
	Values: []string{
		"SIT", "DOWN", "HEEL",
	},
}

// enum CatCommand { JUMP }
var CatCommand = EnumSchema{
	Name:        "CatCommand",
	Description: "Commands that a Cat may know",
	Values:      []string{"JUMP"},
}

// interface Sentient {
//   name: String!
// }
var Sentient = InterfaceSchema{
	Name: "Sentient",
	Fields: Fields{
		"name": NonNullStringType,
	},
}

// interface Pet {
//   name: String!
// }
var Pet = InterfaceSchema{
	Name: "Pet",
	Fields: Fields{
		"name": NonNullStringType,
	},
}

// type Alien implements Sentient {
//   name: String!
//   homePlanet: String
// }
var Alien = ObjectSchema{
	Name:       "Alien",
	Implements: Sentient,
	Fields: ObjectFields{
		"name": ObjectFieldSchema{
			Description: "Name of this Alien",
			Type:        NonNullStringType,
		},
		"homePlanet": ObjectFieldSchema{
			Description: "The name of the planet where this Alien is from",
			Type:        StringType,
		},
	},
}

// type Human implements Sentient {
//   name: String!
// }
var Human = ObjectSchema{
	Name:       "Human",
	Implements: Sentient,
	Fields: ObjectFields{
		"name": ObjectFieldSchema{
			Description: "Name of this Human",
			Type:        NonNullStringType,
		},
	},
}

// type Dog implements Pet {
//   name: String!
//   nickname: String
//   barkVolume: Int
//   doesKnowCommand(dogCommand: DogCommand!): Boolean!
//   isHousetrained(atOtherHomes: Boolean): Boolean!
//   owner: Human
// }
var Dog = ObjectSchema{
	Name:       "Dog",
	Implements: Pet,
	Fields: ObjectFields{
		"name": ObjectFieldSchema{
			Description: "Name of this Dog",
			Type:        NonNullStringType,
		},
		"nickname": ObjectFieldSchema{
			Description: "Nickname of this Dog",
			Type:        StringType,
		},
		"barkVolume": ObjectFieldSchema{
			Description: "How loud this Dog will bark",
			Type:        IntType,
		},
		"doesKnowCommand": ObjectFieldSchema{
			Description: "Function to determine if this Dog knows a given DogCommand",
			Arguments: Arguments{
				"dogCommand": ArgumentSchema{
					Type: DescribeNonNullType("DogCommand"),
				},
			},
			Type: NonNullBooleanType,
		},
		"isHousetrained": ObjectFieldSchema{
			Description: "Function to determine if this Dog is house trained",
			Arguments: Arguments{
				"atOtherHomes": ArgumentSchema{
					Type: BooleanType,
				},
			},
			Type: NonNullBooleanType,
		},
		"owner": ObjectFieldSchema{
			Description: "Owner of this dog",
		},
	},
}

// type Cat implements Pet {
//   name: String!
//   nickname: String
//   doesKnowCommand(catCommand: CatCommand!): Boolean!
//   meowVolume: Int
// }
var Cat = ObjectSchema{
	Name:       "Cat",
	Implements: Pet,
	Fields: ObjectFields{
		"name": ObjectFieldSchema{
			Description: "Name of this Cat",
			Type:        NonNullStringType,
		},
		"nickname": ObjectFieldSchema{
			Description: "Nickname of this Cat",
			Type:        StringType,
		},
		"doesKnowCommand": ObjectFieldSchema{
			Description: "Function to determine if this Cat know a given CatCommand",
			Arguments: Arguments{
				"catCommand": ArgumentSchema{
					Type: NonNullBooleanType,
				},
			},
			Type: NonNullBooleanType,
		},
		"meowVolume": ObjectFieldSchema{
			Description: "How loud this cat meows",
			Type:        IntType,
		},
	},
}

// union CatOrDog = Cat | Dog
var CatOrDog = UnionSchema{
	Name:        "CatOrDog",
	Description: "A type that can either be a Cat or Dog",
	Types:       []ObjectSchema{Cat, Dog},
}

// union DogOrHuman = Dog | Human
var DogOrHuman = UnionSchema{
	Name:        "DogOrHuman",
	Description: "A type that can either be a Dog or Human",
	Types:       []ObjectSchema{Dog, Human},
}

// union HumanOrAlien = Human | Alien
var HumanOrAlien = UnionSchema{
	Name:        "HumanOrAlien",
	Description: "A type that can either be a Human or Alien",
	Types:       []ObjectSchema{Human, Alien},
}

// type QueryRoot {
//   dog: Dog
// }
var QueryRoot = ObjectSchema{
	Name:        "QueryRoot",
	Description: "The query root object for this GraphQL Schema",
	Fields: ObjectFields{
		"dog": ObjectFieldSchema{
			Type: DescribeType("Dog"),
		},
	},
}

// type WithList {
//  names: [String!]!
// }
var WithList = ObjectSchema{
	Name:        "WithList",
	Description: "Test type with field of list type",
	Fields: ObjectFields{
		"names": ObjectFieldSchema{
			// Non-null list of non-null Strings
			Type: DescribeNonNullListType(NonNullStringType),
		},
	},
}
