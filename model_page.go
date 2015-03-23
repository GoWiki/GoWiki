package main

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

type Page struct {
	Current Event
	History History
}

func GetPage(t *bolt.Tx, Name string) (*Page, error) {
	tx := &WikiTx{t}
	page := &Page{}
	if Name == "" {
		Name = "/"
	}
	pagedata := tx.Pages().Names().Get([]byte(Name))
	err := json.Unmarshal(pagedata, page)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (p Page) Save(t *bolt.Tx, Name string) error {
	tx := &WikiTx{t}
	pagedata, err := json.Marshal(p)
	if err != nil {
		return err
	}
	if Name == "" {
		Name = "/"
	}
	b_names := tx.Pages().Names()
	return b_names.Put([]byte(Name), pagedata)
}
