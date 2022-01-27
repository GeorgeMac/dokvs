package boltdb

import (
	"os"
	"testing"

	"github.com/georgemac/dokvs/pkg/kv"
	kvtesting "github.com/georgemac/dokvs/pkg/kv/testing"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
)

func TestBoltDB_KVStore_TestingHarness(t *testing.T) {
	kvtesting.TestHarness(t, func(t *testing.T, data kvtesting.SeedStore) kv.Store {
		db, cleanup := newBoltDB("testing.bolt")
		t.Cleanup(cleanup)

		for _, keyspace := range data.Keyspaces {
			db.Update(func(tx *bolt.Tx) error {
				bkt, err := tx.CreateBucket(keyspace.Name)
				require.NoError(t, err)

				for _, entry := range keyspace.Data {
					require.NoError(t, bkt.Put(entry[0], entry[1]))
				}

				return nil
			})
		}

		return New(db)
	})
}

func newBoltDB(path string) (*bolt.DB, func()) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		panic(err)
	}

	return db, func() {
		db.Close()

		os.Remove(path)
	}
}
