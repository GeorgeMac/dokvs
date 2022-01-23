package etcdoc

type Recipe struct {
	ID string
}

type RecipeSchema struct{}

func (RecipeSchema) Collection() []byte { return []byte("recipes") }

func (RecipeSchema) PrimaryKey(r Recipe) []byte {
	return []byte(r.ID)
}

func ExampleCollection() {
	var kv KV
	recipes := NewCollection(kv, Schema[Recipe](RecipeSchema{}), Serializer[Recipe](JSONSerializer[Recipe]{}))
	_ = recipes
}
