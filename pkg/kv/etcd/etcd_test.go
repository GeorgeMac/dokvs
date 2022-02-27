package etcd

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/georgemac/dokvs/pkg/kv"
	kvtesting "github.com/georgemac/dokvs/pkg/kv/testing"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/tests/v3/integration"
)

func TestEtcd_KVStore_TestingHarness(t *testing.T) {
	kvtesting.TestHarness(t, func(t *testing.T, data kvtesting.SeedStore) kv.Store {
		db, cleanup := newETCD(t)
		t.Cleanup(cleanup)

		ctx := context.Background()
		for _, keyspace := range data.Keyspaces {
			for _, entry := range keyspace.Data {
				key := strings.Join([]string{string(keyspace.Name), string(entry[0])}, "/")
				db.Put(ctx, key, string(entry[1]))
			}
		}

		return New(db)
	})
}

func newETCD(t *testing.T) (clientv3.KV, func()) {
	t.Helper()

	integration.BeforeTest(t)

	var (
		maxReqBytes = 1.5 * 1024 * 1024
		quota       = int64(int(maxReqBytes*1.2) + 8*os.Getpagesize())
	)

	cluster := integration.NewClusterV3(
		t,
		&integration.ClusterConfig{
			Size:                     1,
			QuotaBackendBytes:        quota,
			ClientMaxCallSendMsgSize: 100 * 1024 * 1024,
		},
	)

	client, err := cluster.ClusterClient()
	if err != nil {
		t.Fatal(err)
	}

	return client.KV, func() {
		cluster.Terminate(t)
	}
}
