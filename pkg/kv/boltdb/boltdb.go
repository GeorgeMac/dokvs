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

var ErrBucketNotExist = errors.New("bucket not exist")

type KV struct {
	db *bolt.DB
}

func Open(path string) (*KV, error) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}

	return &KV{db: db}, nil
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

func (k KeyspaceView) Range(_ context.Context, rng kv.RangeOptions) (items []kv.Item, err error) {
	if rng.Limit < 1 {
		rng.Limit = defaultLimit
	}

	if rng.End == nil {
		if v := k.bucket.Get(rng.Start); v != nil {
			items = []kv.Item{{K: rng.Start, V: v}}
		}

		return
	}

	cursor := k.bucket.Cursor()

	key, value := cursor.First()
	if rng.Start != nil {
		key, value = cursor.Seek(rng.Start)
	}

	noEnd := len(rng.End) == 1 && rng.End[0] == '\x00'
	for ; key != nil && (noEnd || bytes.Compare(key, rng.End) <= 0); key, value = cursor.Next() {
		items = append(items, kv.Item{K: key, V: value})
		if len(items) >= rng.Limit {
			break
		}
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

func (u KeyspaceUpdate) Range(ctx context.Context, opts kv.RangeOptions) ([]kv.Item, error) {
	return KeyspaceView(u).Range(ctx, opts)
}

func (u KeyspaceUpdate) Put(_ context.Context, k, v []byte) error {
	return u.bucket.Put(k, v)
}

func (u KeyspaceUpdate) Delete(_ context.Context, k []byte) error {
	return u.bucket.Delete(k)
}
