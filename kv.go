package etcdoc

import "context"

type KVListPredicate struct {
	Prefix []byte
	Offset []byte
	Limit  int
}

type Item struct {
	K, V []byte
}

type KV interface {
	Fetch(context.Context, []byte) (Item, error)
	List(context.Context, KVListPredicate) ([]Item, error)
	Put(_ context.Context, k, v []byte) error
	Delete(_ context.Context, k []byte) error
}