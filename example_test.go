package etcdoc

import "context"

type ID string

type Recipe struct {
	ID ID
}

type RecipeSchema struct{}

func (RecipeSchema) Collection() []byte { return []byte("recipes") }

func (RecipeSchema) PrimaryKey(r Recipe) []byte {
	return []byte(r.ID)
}

func ExampleCollection() {
	recipes := NewCollection[Recipe, ID](Schema[Recipe](RecipeSchema{}), Serializer[Recipe](JSONSerializer[Recipe]{}))

	ctx := context.Background()
	var update KVUpdate
	_ = recipes.Update(update).Put(ctx, Recipe{ID: "my_recipe"})

	var view KVView
	_, _ = recipes.View(view).Fetch(ctx, ID("my_recipe"))
}
