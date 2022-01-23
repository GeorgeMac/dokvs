package etcdoc

import "context"

type ID string

var (
	schema     Schema[Recipe]     = (RecipeSchema{})
	serializer Serializer[Recipe] = (JSONSerializer[Recipe]{})
)

type Recipe struct {
	ID ID
}

type RecipeSchema struct{}

func (RecipeSchema) Collection() []byte { return []byte("recipes") }

func (RecipeSchema) PrimaryKey(r Recipe) []byte {
	return []byte(r.ID)
}

func ExampleCollection() {
	recipes := NewCollection[Recipe, ID](schema, WithSerializer[Recipe, ID](serializer))

	ctx := context.Background()
	var update KVUpdate
	_ = recipes.Update(update).Put(ctx, Recipe{ID: "my_recipe"})

	var view KVView
	_, _ = recipes.View(view).Fetch(ctx, ID("my_recipe"))
}
