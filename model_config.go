package main

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

type Config struct {
	InitDone    bool
	Theme       string
	FilesLoaded bool
}

func GetConfig(t *bolt.Tx) *Config {
	tx := &WikiTx{t}
	configdata := tx.Config().Get([]byte("Main"))
	c := &Config{}
	json.Unmarshal(configdata, c)
	return c
}

func (c *Config) Save(t *bolt.Tx) {
	tx := &WikiTx{t}
	configdata, _ := json.Marshal(c)
	tx.Config().Put([]byte("Main"), configdata)
}
