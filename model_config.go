package main

import "github.com/boltdb/bolt"

type Config struct {
	SetupToken string
	Theme      string
}

func GetConfig(t *bolt.Tx) *Config {
	tx := TX{t}
	tx.Config()

}
