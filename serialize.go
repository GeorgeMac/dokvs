package dokvs

import "encoding/json"

type JSONSerializer[T any] struct{}

func (s JSONSerializer[T]) Serialize(t T) ([]byte, error) {
	return json.Marshal(t)
}

func (d JSONSerializer[T]) Deserialize(v []byte, t *T) error {
	return json.Unmarshal(v, t)
}
