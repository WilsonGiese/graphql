package graphql

type Schema struct {
}

type Scalar struct {
	Name        string
	Description string
}

type EnumSchema struct {
	Name        string
	Values      []string
	Description string
}

type InputSchema struct {
	Name        string
	Description string
}

type Fields map[string]TypeSchema

type InterfaceSchema struct {
	Name        string
	Description string
	Fields      Fields
}

type UnionSchema struct {
	Name        string
	Description string
	Types       []ObjectSchema
}

type ObjectSchema struct {
	Name        string
	Implements  InterfaceSchema
	Description string
	Fields      ObjectFields
}

type ObjectFieldSchema struct {
	Name        string
	Description string
	Type        TypeSchema
	Arguments
}

type ObjectFields map[string]ObjectFieldSchema

type TypeSchema struct {
	Name    string
	NonNull bool
	List    bool
	SubType *TypeSchema
}

type ArgumentSchema struct {
	Name    string
	Type    TypeSchema
	Default interface{}
}

type Arguments map[string]ArgumentSchema

// Type helpers
func DescribeType(name string) (t TypeSchema) {
	t.Name = name
	return
}

func DescribeNonNullType(name string) (t TypeSchema) {
	t.Name = name
	t.NonNull = true
	return
}

func DescribeListType(subType TypeSchema) (t TypeSchema) {
	t.List = true
	t.SubType = &subType
	return
}

func DescribeNonNullListType(subType TypeSchema) (t TypeSchema) {
	t.List = true
	t.NonNull = true
	t.SubType = &subType
	return
}

// Common Pre-defined Types
var StringType = DescribeType("String")

var NonNullStringType = DescribeNonNullType("String")

var IntType = DescribeType("Int")

var NonNullIntType = DescribeNonNullType("Int")

var FloatType = DescribeType("Float")

var NonNullFloatType = DescribeNonNullType("Float")

var BooleanType = DescribeType("Boolean")

var NonNullBooleanType = DescribeNonNullType("Boolean")

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
