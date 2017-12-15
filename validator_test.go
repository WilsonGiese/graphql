package graphql

import (
	"fmt"
	"strings"
	"testing"

	schema "github.com/WilsonGiese/graphql/schema"
)

var SampleSchema *schema.Schema

func Test(t *testing.T) {
	query := `query TestQuery {
  dog {
    name
    doesKnowCommand(dogCommand: SIT)
  }
}
`
	tokens, lexErr := Tokenize(strings.NewReader(query), true)
	fmt.Println(tokens)
	if lexErr == nil {
		document, parseErr := Parse(tokens)
		fmt.Printf("%+v\n", document)
		if parseErr == nil {
			errors := validate(SampleSchema, &document)
			if len(errors) > 0 {
				panic(fmt.Sprintf("%v", errors))
			}
		} else {
			panic(parseErr)
		}
	} else {
		panic(lexErr)
	}
}

func init() {
	SampleSchema = schema.NewSchema().
		Declare(schema.Enum{
			Description: "Commands that a Dog may know",
			Name:        "DogCommand",
			Values:      schema.Values("SIT", "DOWN", "HEEL"),
		}).
		Declare(schema.Enum{
			Description: "Commands that a Cat may know",
			Name:        "CatCommand",
			Values:      schema.Values("JUMP"),
		}).
		Declare(schema.Interface{
			Description: "Sometimes does the thinky thinky",
			Name:        "Sentient",
			Fields: schema.Fields(
				schema.Field{
					Name: "name",
					Type: schema.NonNullStringType,
				},
			),
		}).
		Declare(schema.Interface{
			Name: "Pet",
			Fields: schema.Fields(
				schema.Field{
					Name: "name",
					Type: schema.NonNullStringType,
				},
			),
		}).
		Declare(schema.Object{
			Description: "Sentient type Alien",
			Name:        "Alien",
			Implements:  schema.Interfaces("Sentient"),
			Fields: schema.Fields(
				schema.Field{
					Description: "Name of this Alien",
					Name:        "name",
					Type:        schema.NonNullStringType,
				},
				schema.Field{
					Description: "The name of the planet where this Alien is from",
					Name:        "homePlanet",
					Type:        schema.StringType,
				},
			),
		}).
		Declare(schema.Object{
			Description: "Sentient type Human",
			Name:        "Human",
			Implements:  schema.Interfaces("Sentient"),
			Fields: schema.Fields(
				schema.Field{
					Description: "Name of this Human",
					Name:        "name",
					Type:        schema.NonNullStringType,
				},
			),
		}).
		Declare(schema.Object{
			Description: "Pet type Dog",
			Name:        "Dog",
			Implements:  schema.Interfaces("Pet"),
			Fields: schema.Fields(
				schema.Field{
					Name:        "name",
					Description: "Name of this Dog",
					Type:        schema.NonNullStringType,
				},
				schema.Field{
					Description: "Nickname of this Dog",
					Name:        "nickname",
					Type:        schema.StringType,
				},
				schema.Field{
					Description: "How loud this Dog will bark",
					Name:        "barkVolume",
					Type:        schema.IntType,
				},
				schema.Field{
					Description: "Function to determine if this Dog knows a given DogCommand",
					Name:        "doesKnowCommand",
					Type:        schema.NonNullBooleanType,
					Arguments: schema.Arguments(
						schema.Argument{
							Name: "dogCommand",
							Type: schema.DescribeNonNullType("DogCommand"),
						},
					),
				},
				schema.Field{
					Description: "Function to determine if this Dog is house trained",
					Name:        "isHousetrained",
					Type:        schema.NonNullBooleanType,
					Arguments: schema.Arguments(
						schema.Argument{
							Name: "atOtherHomes",
							Type: schema.BooleanType,
						},
					),
				},
				schema.Field{
					Description: "Owner of this dog",
					Name:        "owner",
					Type:        schema.DescribeType("Human"),
				},
			),
		}).
		Declare(schema.Object{
			Description: "Pet type Cat",
			Name:        "Cat",
			Implements:  schema.Interfaces("Pet"),
			Fields: schema.Fields(
				schema.Field{
					Description: "Name of this Cat",
					Name:        "name",
					Type:        schema.NonNullStringType,
				},
				schema.Field{
					Description: "Nickname of this Cat",
					Name:        "nickname",
					Type:        schema.StringType,
				},
				schema.Field{
					Description: "Function to determine if this Cat know a given CatCommand",
					Name:        "doesKnowCommand",
					Type:        schema.NonNullBooleanType,
					Arguments: schema.Arguments(
						schema.Argument{
							Name: "catCommand",
							Type: schema.DescribeNonNullType("CatCommand"),
						},
					),
				},
				schema.Field{
					Description: "How loud this cat meows",
					Name:        "meowVolume",
					Type:        schema.IntType,
				},
			),
		}).
		Declare(schema.Union{
			Description: "A type that can either be a Cat or Dog",
			Name:        "CatOrDog",
			Types:       schema.Types("Cat", "Dog"),
		}).
		Declare(schema.Union{
			Description: "A type that can either be a Dog or Human",
			Name:        "DogOrHuman",
			Types:       schema.Types("Dog", "Human"),
		}).
		Declare(schema.Union{
			Description: "A type that can either be a Human or Alien",
			Name:        "HumanOrAlien",
			Types:       schema.Types("Human", "Alien"),
		}).
		Declare(schema.Object{
			Description: "The query root object for this GraphQL Schema",
			Name:        "QueryRoot",
			Fields: schema.Fields(
				schema.Field{
					Name: "dog",
					Type: schema.DescribeType("Dog"),
				},
			),
		}).Build()
}
