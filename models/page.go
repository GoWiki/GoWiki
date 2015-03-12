package models

import (
	"github.com/boltdb/bolt"
)

type Page struct {
	Current Event
	History History
}

func GetPage(tx *bolt.Tx) error {
	p := tx.Bucket(pages)

}
