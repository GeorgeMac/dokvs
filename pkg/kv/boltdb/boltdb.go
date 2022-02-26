package boltdb

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/georgemac/dokvs/pkg/kv"
	bolt "go.etcd.io/bbolt"
)

const defaultLimit = 100

var (
	ErrBucketNotExist = errors.New("bucket not exist")

	_ kv.Store = (*KV)(nil)
)

type KV struct {
	db *bolt.DB
}

func New(db *bolt.DB) *KV {
	return &KV{db: db}
}

func (kv KV) Close() error {
	return kv.db.Close()
}

func (kv KV) View(fn func(kv.View) error) error {
	return kv.db.View(func(tx *bolt.Tx) error {
		return fn(View{tx: tx})
	})
}

type View struct {
	tx *bolt.Tx
}

func (v View) Keyspace(key []byte) (_ kv.KeyspaceView, err error) {
	view := KeyspaceView{}
	if view.bucket = v.tx.Bucket(key); view.bucket == nil {
		err = fmt.Errorf("keyspace %q: ", ErrBucketNotExist)
		return
	}

	return view, nil
}

type KeyspaceView struct {
	bucket *bolt.Bucket
}

func (k KeyspaceView) Get(_ context.Context, opts kv.GetOptions) (items []kv.Item, err error) {
	var berr *kv.BatchError
	items = make([]kv.Item, len(opts.Keys))

	for i := range opts.Keys {
		items[i].K = opts.Keys[i]
		if items[i].V = k.bucket.Get(opts.Keys[i]); items[i].V == nil {
			if berr == nil {
				berr = &kv.BatchError{
					Errors: make([]error, len(opts.Keys)),
				}
			}

			berr.Errors[i] = kv.ErrKeyNotFound
		}
	}

	// only assign berr to err if not-nil to avoid
	// returning a (*BatchError)(nil) as err which
	// will would result in the returned err != nil.
	if berr != nil {
		err = berr
	}

	return
}

func (k KeyspaceView) Range(_ context.Context, opts ...kv.RangeOption) (items []kv.Item, err error) {
	var rng kv.RangeOptions
	for _, opt := range opts {
		opt(&rng)
	}

	if rng.Limit < 1 {
		rng.Limit = defaultLimit
	}

	cursor := k.bucket.Cursor()

	key, value := cursor.First()
	if rng.Start != nil {
		key, value = cursor.Seek(rng.Start)
	}

	for ; key != nil && (rng.End == nil || bytes.Compare(key, rng.End) < 0); key, value = cursor.Next() {
		if len(items) >= rng.Limit {
			break
		}
		items = append(items, kv.Item{K: key, V: value})
	}

	return
}

func (kv KV) Update(fn func(kv.Update) error) error {
	return kv.db.Update(func(tx *bolt.Tx) error {
		return fn(Update{tx: tx})
	})
}

type Update struct {
	tx *bolt.Tx
}

func (u Update) CreateKeyspace(key []byte) error {
	_, err := u.tx.CreateBucket(key)
	return err
}

func (u Update) Keyspace(key []byte) (_ kv.KeyspaceUpdate, err error) {
	update := KeyspaceUpdate{}
	if update.bucket = u.tx.Bucket(key); update.bucket == nil {
		err = fmt.Errorf("keyspace %q: ", ErrBucketNotExist)
		return
	}

	return update, nil
}

type KeyspaceUpdate KeyspaceView

func (u KeyspaceUpdate) Get(ctx context.Context, opts kv.GetOptions) ([]kv.Item, error) {
	return KeyspaceView(u).Get(ctx, opts)
}

func (u KeyspaceUpdate) Range(ctx context.Context, opts ...kv.RangeOption) ([]kv.Item, error) {
	return KeyspaceView(u).Range(ctx, opts...)
}

func (u KeyspaceUpdate) Put(_ context.Context, k, v []byte) error {
	return u.bucket.Put(k, v)
}

func (u KeyspaceUpdate) Delete(_ context.Context, k []byte) error {
	return u.bucket.Delete(k)
}
