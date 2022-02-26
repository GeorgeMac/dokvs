package testing

import (
	"context"
	"testing"

	"github.com/georgemac/dokvs/pkg/kv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type SeedKeyspace struct {
	Name []byte
	Data [][2][]byte
}

type SeedStore struct {
	Keyspaces []SeedKeyspace
}

type StoreFactory func(t *testing.T, data SeedStore) kv.Store

func TestHarness(t *testing.T, fn StoreFactory) {
	for _, test := range []struct {
		name string
		seed SeedStore
		test func(t *testing.T, store kv.Store)
	}{
		{
			name: `Keyspace("one")`,
			seed: SeedStore{
				Keyspaces: []SeedKeyspace{
					{
						Name: []byte("one"),
						Data: [][2][]byte{
							{[]byte("a"), []byte("value_one")},
							{[]byte("b"), []byte("value_two")},
							{[]byte("c"), []byte("value_three")},
						},
					},
				},
			},
			test: func(t *testing.T, store kv.Store) {
				ctx := context.Background()

				store.View(func(view kv.View) error {
					keyspace, err := view.Keyspace([]byte("one"))
					require.NoError(t, err)

					t.Run(`Get(Key("a")) is present`, func(t *testing.T) {
						items, err := keyspace.Get(ctx, kv.Key([]byte("a")))
						require.NoError(t, err)

						expected := []kv.Item{{K: []byte("a"), V: []byte("value_one")}}
						assert.Equal(t, expected, items)
					})

					t.Run(`Get(Key("d")) returns not found`, func(t *testing.T) {
						_, err := keyspace.Get(ctx, kv.Key([]byte("d")))
						require.Equal(
							t,
							&kv.BatchError{
								Errors: []error{kv.ErrKeyNotFound},
							},
							err,
						)
					})

					t.Run(`Get(Batch("a", "b")) returns ["a", "b"]`, func(t *testing.T) {
						items, err := keyspace.Get(ctx, kv.Batch([]byte("a"), []byte("b")))
						require.NoError(t, err)

						expected := []kv.Item{
							{K: []byte("a"), V: []byte("value_one")},
							{K: []byte("b"), V: []byte("value_two")},
						}
						assert.Equal(t, expected, items)
					})

					t.Run(`Get(Batch("a", "d")) returns ["a", <not found>]`, func(t *testing.T) {
						items, err := keyspace.Get(ctx, kv.Batch([]byte("a"), []byte("d")))

						expected := []kv.Item{
							{K: []byte("a"), V: []byte("value_one")},
							{K: []byte("d")}, // item not present
						}
						assert.Equal(t, expected, items)
						require.Equal(
							t,
							&kv.BatchError{
								Errors: []error{
									nil,
									kv.ErrKeyNotFound, // error reflects missing item d
								},
							},
							err,
						)
					})

					t.Run(`Range(["a", "c")) returns ["a", "b"]`, func(t *testing.T) {
						items, err := keyspace.Range(
							ctx,
							kv.Start([]byte("a")),
							kv.End([]byte("c")),
						)
						require.NoError(t, err)

						expected := []kv.Item{
							{K: []byte("a"), V: []byte("value_one")},
							{K: []byte("b"), V: []byte("value_two")},
						}
						assert.Equal(t, expected, items)
					})

					t.Run(`Range(["a", *)) returns ["a", "b", "c"]`, func(t *testing.T) {
						items, err := keyspace.Range(
							ctx,
							kv.Start([]byte("a")),
						)
						require.NoError(t, err)

						expected := []kv.Item{
							{K: []byte("a"), V: []byte("value_one")},
							{K: []byte("b"), V: []byte("value_two")},
							{K: []byte("c"), V: []byte("value_three")},
						}
						assert.Equal(t, expected, items)
					})

					t.Run(`Range(["a", *), Limit(2)) returns ["a", "b"]`, func(t *testing.T) {
						items, err := keyspace.Range(
							ctx,
							kv.Start([]byte("a")),
							kv.Limit(2),
						)
						require.NoError(t, err)

						expected := []kv.Item{
							{K: []byte("a"), V: []byte("value_one")},
							{K: []byte("b"), V: []byte("value_two")},
						}
						assert.Equal(t, expected, items)
					})

					return nil
				})
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			store := fn(t, test.seed)

			test.test(t, store)
		})
	}
}
