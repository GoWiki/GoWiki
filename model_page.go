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
	page := &Page{}
	if Name == "" {
		Name = "/"
	}
	pagedata := p.Bucket(names).Get([]byte(Name))
	err := json.Unmarshal(pagedata, page)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (p Page) Save(tx *bolt.Tx, Name string) error {
	pagedata, err := json.Marshal(p)
	if err != nil {
		return err
	}
	if Name == "" {
		Name = "/"
	}
	b_names := tx.Bucket(pages).Bucket(names)
	return b_names.Put([]byte(Name), pagedata)
}
