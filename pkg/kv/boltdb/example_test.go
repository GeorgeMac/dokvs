package boltdb

import (
	"context"
	"fmt"

	"github.com/georgemac/dokvs"
	"github.com/georgemac/dokvs/pkg/kv"
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

	db, cleanup := newBoltDB("example.bolt")
	defer cleanup()

    store := New(db)

	if err := store.Update(func(update kv.Update) error {
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

		recipe, err := recipes.Fetch(ctx, ID("my_recipe"))
		if err != nil {
			return err
		}

		fmt.Printf("%#v\n", recipe)

		if err := recipes.Put(ctx, Recipe{ID: "second_recipe"}); err != nil {
			return err
		}

		allRecipes, err := recipes.List(ctx, dokvs.ListPredicate{})
		if err != nil {
			return err
		}

		fmt.Printf("%#v\n", allRecipes)

		return nil
	}); err != nil {
		panic(err)
	}

	// OUTPUT: boltdb.Recipe{ID:"my_recipe"}
	// []boltdb.Recipe{boltdb.Recipe{ID:"my_recipe"}, boltdb.Recipe{ID:"second_recipe"}}
}
