package main

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

var (
	bn_pages   []byte = []byte("pages")
	bn_history []byte = []byte("history")
	bn_names   []byte = []byte("names")
	bn_data    []byte = []byte("data")
	bn_users   []byte = []byte("users")
)

func SetupBuckets(tx *bolt.Tx) {
	pages, _ := tx.CreateBucketIfNotExists(bn_pages)
	pages.CreateBucketIfNotExists(bn_history)
	pages.CreateBucketIfNotExists(bn_names)
	pages.CreateBucketIfNotExists(bn_data)
	tx.CreateBucketIfNotExists(bn_users)
}

func NextKey(b *bolt.Bucket) []byte {
	i, _ := b.NextSequence()
	key, _ := json.Marshal(i)
	return key
}
