package schema

import (
	"testing"
)

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
	Types:       []string{"Cat", "Dog"},
}

// union DogOrHuman = Dog | Human
var DogOrHuman = Union{
	Name:        "DogOrHuman",
	Description: "A type that can either be a Dog or Human",
	Types:       []string{"Dog", "Human"},
}

// union HumanOrAlien = Human | Alien
var HumanOrAlien = Union{
	Name:        "HumanOrAlien",
	Description: "A type that can either be a Human or Alien",
	Types:       []string{"Human", "Alien"},
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
func Test(t *testing.T) {
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
		// Declare(Scalar{
		// 	Name:        "space ",
		// 	Description: "Represents a datetime with an ISO8601 format",
		// }).
		// Declare(Input{
		// 	Name: "NewDog",
		// 	Fields: Fields(
		// 		Field{
		// 			Name: "Friends",
		// 			Type: StringType,
		// 			Arguments: Arguments(
		// 				Argument{
		// 					Name: "Arg1",
		// 					Type: DescribeListType(Type{Name: "abc"}),
		// 				},
		// 			),
		// 		},
		// 		Field{
		// 			Name: "BarkVolume",
		// 			Type: IntType,
		// 		},
		// 	),
		// }).
		Declare(Enum{
			Name:        "DogCommand",
			Description: "Commands that a Dog may know",
			Values:      []string{"SIT", "DOWN", "HEEL"},
		}).
		// Declare(Enum{
		// 	Name:        "DogCommand",
		// 	Description: "Commands that a Dog may know",
		// 	Values:      []string{"SIT", "DOWN", "HEEL"},
		// }).
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
			//Types:       []string{"Cat", "Dog"},
		}).
		Declare(Union{
			Name:        "DogOrHuman",
			Description: "A type that can either be a Dog or Human",
			Types:       []string{"Dog", "Human"},
		}).
		Declare(Union{
			Name:        "HumanOrAlien",
			Description: "A type that can either be a Human or Alien",
			Types:       []string{"Human", "Alien"},
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
