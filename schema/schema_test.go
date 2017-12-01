package schema

import (
	"testing"
)

func Test(t *testing.T) {
	NewSchema().
		Declare(Scalar{
			Name:        "space",
			Description: "Represents a datetime with an ISO8601 format",
		}).
		Declare(Input{
			Name: "NewDog",
			Fields: Fields(
				Field{
					Name: "f1",
					Type: StringType,
				},
				Field{
					Name: "BarkVolume",
					Type: IntType,
				},
			),
		}).
		Declare(Enum{
			Name:        "DogCommand",
			Description: "Commands that a Dog may know",
			Values:      Values("SIT", "DOWN", "HEEL"),
		}).
		// Declare(Enum{
		// 	Name:        "DogCommand",
		// 	Description: "Commands that a Dog may know",
		// 	Values:      []string{"SIT", "DOWN", "HEEL"},
		// }).
		Declare(Enum{
			Name:        "CatCommand",
			Description: "Commands that a Cat may know",
			Values:      Values("JUMP"),
		}).
		Declare(Interface{
			Name:        "Sentient",
			Description: "Sometimes does the thinky thinky",
			Fields: Fields(
				Field{
					Name: "name",
					Type: NonNullStringType,
					Arguments: Arguments(
						Argument{
							Name: "language",
							Type: StringType,
						},
					),
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
					Arguments: Arguments(
						Argument{
							Name: "language",
							Type: StringType,
						},
						Argument{
							Name: "x",
							Type: StringType,
						},
					),
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
					Arguments: Arguments(
						Argument{
							Name: "language",
							Type: StringType,
						},
					),
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
			Types:       Types("Cat", "Dog"),
		}).
		Declare(Union{
			Name:        "DogOrHuman",
			Description: "A type that can either be a Dog or Human",
			Types:       Types("Dog", "Human"),
		}).
		Declare(Union{
			Name:        "HumanOrAlien",
			Description: "A type that can either be a Human or Alien",
			Types:       Types("Human", "Alien"),
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
