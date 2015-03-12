package models

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

type Page struct {
	Current Event
	History History
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
