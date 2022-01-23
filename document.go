package etcdoc

import (
	"context"
	"encoding/json"
	"errors"
)

var ErrNotFound = errors.New("document not found")

type Serializer[T any] interface {
	Serialize(T) ([]byte, error)
	Deserialize([]byte, *T) error
}

type JSONSerializer[T any] struct{}

func (s JSONSerializer[T]) Serialize(t T) ([]byte, error) {
	return json.Marshal(t)
}

func (d JSONSerializer[T]) Deserialize(v []byte, t *T) error {
	return json.Unmarshal(v, t)
}

type Schema[D any] interface {
	Collection() []byte
    PrimaryKey(D) []byte
}

type Collection[D any] struct {
	kv         KV
	schema     Schema[D]
	serializer Serializer[D]
}

func NewCollection[D any](kv KV, schema Schema[D], serializer Serializer[D]) Collection[D] {
    return Collection[D]{kv: kv, schema: schema, serializer: serializer}
}

func (c *Collection[D]) Fetch(ctx context.Context, key []byte) (d D, err error) {
	item, err := c.kv.Fetch(ctx, key)
	if err != nil {
		return d, err
	}

	err = c.serializer.Deserialize(item.V, &d)

	return
}

type ListPredicate struct {
	Offset []byte
	Limit  int
}

func (c *Collection[D]) List(context.Context, ...ListPredicate) (d D, err error) {
	return
}

func (c *Collection[D]) Put(ctx context.Context, doc D) error {
	v, err := c.serializer.Serialize(doc)
	if err != nil {
		return err
	}

	return c.kv.Put(ctx, c.schema.PrimaryKey(doc), v)
}

func (c *Collection[D]) Delete(ctx context.Context, doc D) error {
	return c.kv.Delete(ctx, c.schema.PrimaryKey(doc))
}
