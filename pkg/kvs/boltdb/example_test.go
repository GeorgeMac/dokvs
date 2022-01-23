package boltdb

import (
	"context"
	"fmt"

	"github.com/georgemac/dokvs"
)

type ID string

var (
	schema = dokvs.NewSchema("recipes", func(r Recipe) []byte {
		return []byte(r.ID)
	})

	serializer dokvs.Serializer[Recipe] = (dokvs.JSONSerializer[Recipe]{})
)

type Recipe struct {
	ID ID
}

func Example_Bolt_Collection() {
	recipes := dokvs.NewCollection[Recipe, ID](schema, dokvs.WithSerializer[Recipe, ID](serializer))

	ctx := context.Background()

	kv, err := Open("example.bolt")
	if err != nil {
		panic(err)
	}

	if err := kv.Update(func(update dokvs.Update) error {
		if err := recipes.Init(update); err != nil {
			return err
		}

		recipes, err := recipes.Update(update)
		if err != nil {
			return err
		}

		if err := recipes.Put(ctx, Recipe{ID: "my_recipe"}); err != nil {
			return err
		}

		items, err := recipes.Fetch(ctx, ID("my_recipe"))
		if err != nil {
			return err
		}

		fmt.Println(items)

		return nil
	}); err != nil {
		panic(err)
	}

	// OUTPUT: {my_recipe}
}
