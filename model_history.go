package main

import (
	"time"

	"github.com/boltdb/bolt"
)

type Event struct {
	DataID   []byte
	Author   []byte
	IP       string
	DateTime time.Time
}

func (p Event) GetData(tx *bolt.Tx) []byte {
	b_data := tx.Bucket(bn_pages).Bucket(bn_data)
	return b_data.Get(p.DataID)
}

func SaveData(tx *bolt.Tx, pagedata []byte) ([]byte, error) {
	b_data := tx.Bucket(bn_pages).Bucket(bn_data)
	key := NextKey(b_data)
	b_data.Put(key, pagedata)
	return key, nil
}

type History struct {
	Events []*Event
}
