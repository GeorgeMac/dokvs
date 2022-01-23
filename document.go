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

type AnyBytes interface {
	~[]byte | ~string
}

type Schema[D any] interface {
	Collection() []byte
	PrimaryKey(D) []byte
}

type Collection[D any, K AnyBytes] struct {
	schema     Schema[D]
	serializer Serializer[D]
}

func NewCollection[D any, K AnyBytes](schema Schema[D], serializer Serializer[D]) Collection[D, K] {
	return Collection[D, K]{schema: schema, serializer: serializer}
}

func (c Collection[D, K]) View(view KVView) CollectionView[D, K] {
	return CollectionView[D, K]{Collection: c, view: view}
}

func (c Collection[D, K]) Update(update KVUpdate) CollectionUpdate[D, K] {
	return CollectionUpdate[D, K]{CollectionView: c.View(update), update: update}
}

type CollectionView[D any, K AnyBytes] struct {
	Collection[D, K]

	view KVView
}

func (c CollectionView[D, K]) Fetch(ctx context.Context, key K) (d D, err error) {
	item, err := c.view.Fetch(ctx, []byte(key))
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

func (c CollectionView[D, K]) List(context.Context, ...ListPredicate) (d D, err error) {
	err = errors.New("TODO: implement")
	return
}

type CollectionUpdate[D any, K AnyBytes] struct {
	CollectionView[D, K]

	update KVUpdate
}

func (c CollectionUpdate[D, K]) Fetch(ctx context.Context, key K) (d D, err error) {
	return c.CollectionView.Fetch(ctx, key)
}

func (c CollectionUpdate[D, K]) List(ctx context.Context, pred ...ListPredicate) (d D, err error) {
	return c.CollectionView.List(ctx, pred...)
}

func (c CollectionUpdate[D, K]) Put(ctx context.Context, doc D) error {
	v, err := c.serializer.Serialize(doc)
	if err != nil {
		return err
	}

	return c.update.Put(ctx, c.schema.PrimaryKey(doc), v)
}

func (c CollectionUpdate[D, K]) Delete(ctx context.Context, doc D) error {
	return c.update.Delete(ctx, c.schema.PrimaryKey(doc))
}
