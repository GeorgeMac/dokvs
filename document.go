package etcdoc

import (
	"context"
	"errors"
    "encoding/json"
)

var ErrNotFound = errors.New("document not found")

type Serializer[T any] interface {
	Serialize(T) ([]byte, error)
    Deserialize([]byte, *T) error
}

type JSONSerializer[T any] struct {}

func (s JSONSerializer[T]) Serialize(t T) ([]byte, error) {
    return json.Marshal(t)
}

func (d JSONSerializer[T]) Deserialize(v []byte, t *T) error {
    return json.Unmarshal(v, t)
}

type Document interface {
	PrimaryKey() []byte
}

type Collection[D Document] struct {
	// client clientv3.KV
	kv         KV
	serializer Serializer[D]
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

	return c.kv.Put(ctx, doc.PrimaryKey(), v)
}

func (c *Collection[D]) Delete(ctx context.Context, doc D) error {
	return c.kv.Delete(ctx, doc.PrimaryKey())
}
