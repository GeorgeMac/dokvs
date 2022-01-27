package kv

import (
	"context"
	"errors"
)

// ErrKeyNotFound is returned when a key is not found.
var ErrKeyNotFound = errors.New("key not found")

// Item is a single entry in a Key/Value keyspace in a Key/Value store.
type Item struct {
	K, V []byte
}

// Store represents a Key/Value store which support both read-only (View)
// and read-write (Update) transactional operations.
type Store interface {
	View(func(View) error) error
	Update(func(Update) error) error
}

// View is a read-only transactional view of a KV store.
// It can be used to execute read-operations across
// one or more of keyspaces.
type View interface {
	Keyspace([]byte) (KeyspaceView, error)
}

// RangeOptions contains a set of Range operation predicates.
// RangeOptions is inspired by ETCD's RangeRequest protobuf structure.
//
// If End it omitted, then Range behaves like a Get.
// It returns 0 or 1 keys for an exact match on the provided Start byte-slice.
//
// Else it returns the set of keys in the range [Start, End).
// If end contains the single null-byte, then it returns the range [Start, *).
// Where * represents no bound on the end of the range (all keys from Start).
type RangeOptions struct {
	// Start is the initial location to begin the range operation.
	// If Start == nil, then the range starts at the beginning of the keyspace.
	// If Start != nil, then the range starts at the first key in the keyspace,
	// which is >= Start.
	Start []byte
	// End is the non-inclusive end of the range.
	// If End == nil, then Range only returns a single item if Start is present in the keyspace (effectively a Get(key)).
	// Else If End == []byte{\x00}, then the range is to the end of the entire keyspace [Start, *).
	// Else the range [Start, End) a.k.a Start >= keys < End.
	End []byte
	// Limit is a limit on the number of returned keys.
	// Limit == 0 means no limit.
	// Limit > 0 ensures at-most Limit items are returned.
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
