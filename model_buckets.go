package main

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

var (
	pages   []byte = []byte("pages")
	history []byte = []byte("history")
	names   []byte = []byte("names")
	data    []byte = []byte("data")
)

func SetupBuckets(tx *bolt.Tx) {
	pages, _ := tx.CreateBucketIfNotExists(pages)
	pages.CreateBucketIfNotExists(history)
	pages.CreateBucketIfNotExists(names)
	pages.CreateBucketIfNotExists(data)
}

func NextKey(b *bolt.Bucket) []byte {
	i, _ := b.NextSequence()
	key, _ := json.Marshal(i)
	return key
}
