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
	var update KVUpdate
	_ = recipes.Update(update).Put(ctx, Recipe{ID: "my_recipe"})

	var view KVView
	_, _ = recipes.View(view).Fetch(ctx, ID("my_recipe"))
}
