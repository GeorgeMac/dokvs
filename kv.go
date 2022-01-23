package dokvs

import "context"

type RangeOptions struct {
	Start []byte
	End   []byte
	Limit int
}

type Item struct {
	K, V []byte
}

type KV interface {
	View(func(View) error) error
	Update(func(Update) error) error
}

type View interface {
	Keyspace([]byte) (KeyspaceView, error)
}

type KeyspaceView interface {
	Range(context.Context, RangeOptions) ([]Item, error)
}

type Update interface {
	CreateKeyspace([]byte) error
	Keyspace([]byte) (KeyspaceUpdate, error)
}

type KeyspaceUpdate interface {
	KeyspaceView

	Put(_ context.Context, k, v []byte) error
	Delete(_ context.Context, k []byte) error
}
