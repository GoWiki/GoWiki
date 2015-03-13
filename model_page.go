package main

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

type Page struct {
	Current Event
	History []byte
}

func GetPage(tx *bolt.Tx, Name string) (*Page, error) {
	p := tx.Bucket(pages)
	var page *Page
	data := p.Bucket(names).Get([]byte(Name))
	err := json.Unmarshal(data, page)

	if err != nil {
		return nil, err
	}
	return page, nil
}

func (p Page) Save(tx *bolt.Tx, Name string) error {
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	b_names := tx.Bucket(pages).Bucket(names)

	b_names.Put([]byte(Name), data)
	return nil
}
