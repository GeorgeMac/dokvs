package dokvs

import "context"

type ID string

var (
	schema = NewSchema("recipes", func(r Recipe) []byte {
		return []byte(r.ID)
	})

	serializer Serializer[Recipe] = (JSONSerializer[Recipe]{})
)

type Recipe struct {
	ID ID
}

func ExampleCollection() {
	recipes := NewCollection[Recipe, ID](schema, WithSerializer[Recipe, ID](serializer))

	ctx := context.Background()

	var update Update
	recipesUpdate, _ := recipes.Update(update)
	_ = recipesUpdate.Put(ctx, Recipe{ID: "my_recipe"})

	var view View
	recipesView, _ := recipes.View(view)
	_, _ = recipesView.Fetch(ctx, ID("my_recipe"))
}
