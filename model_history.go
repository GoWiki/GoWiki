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

func (p Event) GetData(t *bolt.Tx) []byte {
	tx := &TX{t}
	b_data := tx.Pages().Data()
	return b_data.Get(p.DataID)
}

func SaveData(t *bolt.Tx, pagedata []byte) ([]byte, error) {
	tx := &TX{t}
	b_data := tx.Pages().Data()
	key := NextKey(b_data)
	b_data.Put(key, pagedata)
	return key, nil
}

type History struct {
	Events []*Event
}

func (h *History) AddEvent(e Event) {
	h.Events = append(h.Events, &e)
}
