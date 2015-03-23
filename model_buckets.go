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
	bn_emails  []byte = []byte("emails")
	bn_config  []byte = []byte("config")
	bn_themes  []byte = []byte("themes")
)

func SetupBuckets(tx *bolt.Tx) {
	pages, _ := tx.CreateBucketIfNotExists(bn_pages)
	pages.CreateBucketIfNotExists(bn_history)
	pages.CreateBucketIfNotExists(bn_names)
	pages.CreateBucketIfNotExists(bn_data)
	users, _ := tx.CreateBucketIfNotExists(bn_users)
	users.CreateBucketIfNotExists(bn_data)
	users.CreateBucketIfNotExists(bn_names)
	users.CreateBucketIfNotExists(bn_emails)
	tx.CreateBucketIfNotExists(bn_config)
	tx.CreateBucketIfNotExists(bn_themes)
}

type WikiTx struct {
	*bolt.Tx
}
type Pages struct {
	*bolt.Bucket
}
type Users struct {
	*bolt.Bucket
}

func (tx *WikiTx) Users() *Users {
	return &Users{tx.Bucket(bn_users)}
}

func (tx *WikiTx) Pages() *Pages {
	return &Pages{tx.Bucket(bn_pages)}
}

func (tx *WikiTx) Config() *bolt.Bucket {
	return tx.Bucket(bn_config)
}

func (tx *WikiTx) Themes() *bolt.Bucket {
	return tx.Bucket(bn_themes)
}

func (p *Pages) Names() *bolt.Bucket {
	return p.Bucket.Bucket(bn_names)
}

func (p *Pages) History() *bolt.Bucket {
	return p.Bucket.Bucket(bn_history)
}

func (p *Pages) Data() *bolt.Bucket {
	return p.Bucket.Bucket(bn_data)
}

func (u *Users) Names() *bolt.Bucket {
	return u.Bucket.Bucket(bn_names)
}

func (u *Users) Emails() *bolt.Bucket {
	return u.Bucket.Bucket(bn_emails)
}

func (u *Users) Data() *bolt.Bucket {
	return u.Bucket.Bucket(bn_data)
}

func NextKey(b *bolt.Bucket) []byte {
	i, _ := b.NextSequence()
	key, _ := json.Marshal(i)
	return key
}
