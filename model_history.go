package main

import (
	"time"

	"github.com/boltdb/bolt"
)

type Event struct {
	DataID   []byte
	AuthorID []byte
	IP       string
	DateTime time.Time
	Author   *SafeUser `json:"-"`
}

func (p Event) GetData(t *bolt.Tx) []byte {
	tx := &WikiTx{t}
	b_data := tx.Pages().Data()
	return b_data.Get(p.DataID)
}

func SaveData(t *bolt.Tx, pagedata []byte) ([]byte, error) {
	tx := &WikiTx{t}
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

func (h *History) LoadUsers(t *bolt.Tx) {
	for _, v := range h.Events {
		v.Author = GetSafeUserByID(t, v.AuthorID)
		if v.Author == nil {
			v.Author = &SafeUser{Name: "anonymous"}
		}
	}
}
