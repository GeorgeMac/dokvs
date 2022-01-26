dokvs - Documents Over KV Stores
-------------------------------

This project aims to provide document-store capabilities, built over key-value persisted backends.
It is built in Go using the new generics features to expose a friendly, type-safe API.

This is mostly a personal exercise for getting used to the new generics features of Go.
By means of applying them to a problem I am familiar with.

## Requirements

1. Go at `>=1.18beta1`

## Backends

- [x] [BoltDB](./pkg/kvs/boltdb)
- [ ] Etcd

## Inspirations

`dokvs` is heavily inspired by my day to day work @InfluxData.

However, it is also inspired by lots of projects in the world of Go.
Including but not limited to:

- [InfluxDBs KV Store Abstractions](https://github.com/influxdata/influxdb/blob/master/kv/store.go#L27-L60)
- [ETCD](https://github.com/etcd-io/etcd)
- [BoltDB](https://github.com/etcd-io/bbolt)
- [Kubernetes client-go](https://github.com/kubernetes/client-go/)
