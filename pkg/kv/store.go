package kv

import (
	"context"
	"errors"
	"fmt"
)

// ErrKeyNotFound is returned when a key is not found.
var ErrKeyNotFound = errors.New("key not found")

// BatchError is a struct containing a slice of errors which also
// implements the error interface.
// It is returned when any item in a batch requested via Get
// results in an error.
// The errors are indexed based on the key requested in the
// GetOptions supplied to Get.
type BatchError struct {
	Errors []error
}

// Error returns a string representations of the BatchError
func (b BatchError) Error() string {
	return fmt.Sprintf("%v", b.Errors)
}

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

// GetOptions defines a set of keys which refer to a set of items in a keyspace.
//
// It is used in a call to KeyspaceView.Get to fetch a batch of items.
type GetOptions struct {
	Keys [][]byte
}

// Key returns a GetOptions which refers to a single key.
func Key(key []byte) GetOptions {
	return Batch(key)
}

// Batch returns a GetOptions which refers to all keys in the set.
func Batch(keys ...[]byte) GetOptions {
	return GetOptions{keys}
}

// RangeOption is a function which configures a RangeOptions
type RangeOption func(*RangeOptions)

// Start configures the inclusive beginning of a Range call.
// See RangeOptions{}.
func Start(key []byte) RangeOption {
	return func(o *RangeOptions) {
		o.Start = key
	}
}

// End configures the exclusive end of a Range call.
// See RangeOptions{}.
func End(key []byte) RangeOption {
	return func(o *RangeOptions) {
		o.End = key
	}
}

// Limit configures the maximum number of items to return in a Range call.
// See RangeOptions{}.
func Limit(n int) RangeOption {
	return func(o *RangeOptions) {
		o.Limit = n
	}
}

// RangeOptions is used when requesting a Range of items from a key/value store.
//
// It us used in a call to KeyspaceView.Range to fetch a sequence of items.
// It defines the start, end and a limit on the number of items returned.
// The zero-value of range options represents a request for the entire keyspace.
type RangeOptions struct {
	// Start defines the inclusive beginning of the range requested in the keyspace.
	// If Start == nil, this refers to the beginning of the keyspace.
	// If Start != nil, this refers to items with key >= start.
	Start []byte
	// End defines the exclusive end of the range requested in the keyspace.
	// If End == nil, this refers to all keys until the end of the keyspace.
	// If End != nil, this refers to items with key < end.
	End []byte
	// Limit is the maximum number of items to return.
	Limit int
}

// KeyspaceView is a read-only client for accessing ranges of a single keyspace
// in a Key/Value store.
type KeyspaceView interface {
	Get(context.Context, GetOptions) ([]Item, error)
	Range(context.Context, ...RangeOption) ([]Item, error)
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
