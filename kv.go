package dokvs

import "context"

// Item is a single entry in a Key/Value Keyspace in a Key/Value Store.
type Item struct {
	K, V []byte
}

// KV represents a KV store which support both read-only (View)
// and writable (Update) transactional operations.
type KV interface {
	View(func(View) error) error
	Update(func(Update) error) error
}

// View is a read-only transactional view of a KV store.
// It can be used to execute read-operations across
// one or more of keyspaces.
type View interface {
	Keyspace([]byte) (KeyspaceView, error)
}

// RangeOptions contains a set of predicates to be supplied
// on Range operation.
// RangeOptions is inspired by ETCD's RangeRequest protobuf structure.
// A range is defined in the form of [Start, End).
type RangeOptions struct {
	// Start is the initial location to begin the range operation.
	// If Start == nil, then the range starts at the beginning of the keyspace.
	// If Start != nil, then the range starts at the first key in the keyspace,
	// which is >= Start. e.g.
	Start []byte
	// End is the non-inclusive end of the range.
	// If End == nil, then the range is to the end of the entire keyspace.
	// If End != nul, then the range [Start, End) a.k.a Start >= keys < End.
	End []byte
	// Limit is a limit on the number of returned keys.
	Limit int
}

// KeyspaceView is a read-only client for accessing ranges of a single keyspace
// in a Key/Value store.
type KeyspaceView interface {
	Range(context.Context, RangeOptions) ([]Item, error)
}

// Update is a read-write transaction of a KV store.
// It supports creating new keyspaces as well as obtaining
// mutable KeyspaceUpdate; used to perform keyspace updates.
type Update interface {
	CreateKeyspace([]byte) error
	Keyspace([]byte) (KeyspaceUpdate, error)
}

// KeyspaceUpdate is a read-write interface across a particular keyspace, which can:
// Range across the keyspace.
// Put into the keyspace.
// Delete from the keyspace.
type KeyspaceUpdate interface {
	KeyspaceView

	Put(_ context.Context, k, v []byte) error
	Delete(_ context.Context, k []byte) error
}
