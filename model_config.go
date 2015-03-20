package main

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

type Config struct {
	InitDone bool
	Theme    string
}

func GetConfig(t *bolt.Tx) *Config {
	tx := TX{t}
	configdata := tx.Config().Get([]byte("Main"))
	c := &Config{}
	json.Unmarshal(configdata, c)
	return c
}

func (c *Config) Save(t *bolt.Tx) {
	tx := &TX{t}
	configdata, _ := json.Marshal(c)
	tx.Config().Put([]byte("Main"), configdata)
}
