package etcd

import (
	"bytes"
	"context"
	"strings"

	"github.com/georgemac/dokvs/pkg/kv"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var _ kv.Store = (*KV)(nil)

type KV struct {
	kv clientv3.KV
}

func New(kv clientv3.KV) *KV {
	return &KV{kv: kv}
}

func (kv KV) View(fn func(kv.View) error) error {
	return fn(View{kv: kv.kv})
}

type View struct {
	kv clientv3.KV
}

func (v View) Keyspace(key []byte) (_ kv.KeyspaceView, err error) {
	return KeyspaceView{
		kv:     v.kv,
		prefix: key,
	}, nil
}

type KeyspaceView struct {
	kv     clientv3.KV
	prefix []byte
}

func (k KeyspaceView) key(v []byte) string {
	return strings.Join([]string{string(k.prefix), string(v)}, "/")
}

func (k KeyspaceView) Get(ctx context.Context, opts kv.GetOptions) (items []kv.Item, err error) {
	getOps := make([]clientv3.Op, len(opts.Keys))
	for i := range opts.Keys {
		getOps[i] = clientv3.OpGet(k.key(opts.Keys[i]))
	}

	resp, err := k.kv.Txn(ctx).
		Then(getOps...).
		Commit()
	if err != nil {
		return nil, err
	}

	var berr *kv.BatchError
	appendError := func(i int, err error) {
		if berr == nil {
			berr = &kv.BatchError{
				Errors: make([]error, len(opts.Keys)),
			}
		}

		berr.Errors[i] = err
	}

	items = make([]kv.Item, len(resp.Responses))
	for i, op := range resp.Responses {
		items[i].K = opts.Keys[i]

		rng := op.GetResponseRange()
		if rng == nil {
			appendError(i, kv.ErrKeyNotFound)
			continue
		}

		if len(rng.Kvs) < 1 {
			appendError(i, kv.ErrKeyNotFound)
			continue
		}

		items[i].V = rng.Kvs[0].Value
	}

	if berr != nil {
		err = berr
	}

	return
}

func (k KeyspaceView) prefixKey(dst *[]byte, src []byte) {
	if src == nil {
		*dst = make([]byte, len(k.prefix)+1)
		copy(*dst, k.prefix[:])
		(*dst)[len(k.prefix)] = '/'
	} else {
		*dst = make([]byte, len(k.prefix)+len(src)+1)
		copy(*dst, k.prefix[:])
		(*dst)[len(k.prefix)] = '/'
		copy((*dst)[len(k.prefix)+1:], src[:])
	}
}

func (k KeyspaceView) Range(ctx context.Context, opts ...kv.RangeOption) (items []kv.Item, err error) {
	var rng kv.RangeOptions
	for _, opt := range opts {
		opt(&rng)
	}

	var start []byte
	k.prefixKey(&start, rng.Start)

	var end []byte
	k.prefixKey(&end, rng.End)

	var rngEnd string
	if rng.End == nil {
		rngEnd = clientv3.GetPrefixRangeEnd(string(end))
	} else {
		rngEnd = string(end)
	}

	resp, err := k.kv.Get(
		ctx,
		string(start),
		clientv3.WithRange(rngEnd),
		clientv3.WithLimit(int64(rng.Limit)),
	)
	if err != nil {
		return nil, err
	}

	items = make([]kv.Item, len(resp.Kvs))
	for i := range resp.Kvs {
		items[i].K = bytes.TrimPrefix(resp.Kvs[i].Key, append(k.prefix, '/'))
		items[i].V = resp.Kvs[i].Value
	}

	return
}

func (kv KV) Update(fn func(kv.Update) error) error {
	return fn(Update{kv: kv.kv})
}

type Update struct {
	kv clientv3.KV
}

// CreateKeyspace is a noop for etcd.
func (u Update) CreateKeyspace(key []byte) error {
	return nil
}

func (u Update) Keyspace(key []byte) (_ kv.KeyspaceUpdate, err error) {
	return KeyspaceUpdate{
		kv:     u.kv,
		prefix: key,
	}, nil
}

type KeyspaceUpdate KeyspaceView

func (u KeyspaceUpdate) Get(ctx context.Context, opts kv.GetOptions) ([]kv.Item, error) {
	return KeyspaceView(u).Get(ctx, opts)
}

func (u KeyspaceUpdate) Range(ctx context.Context, opts ...kv.RangeOption) ([]kv.Item, error) {
	return KeyspaceView(u).Range(ctx, opts...)
}

func (u KeyspaceUpdate) Put(ctx context.Context, k, v []byte) error {
	_, err := u.kv.Put(ctx, string(k), string(v))
	return err
}

func (u KeyspaceUpdate) Delete(ctx context.Context, k []byte) error {
	_, err := u.kv.Delete(ctx, string(k))
	return err
}
