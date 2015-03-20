package main

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       []byte
	Name     string
	Password []byte
}

func GetUser(t *bolt.Tx, name string) *User {
	tx := TX{t}
	userid := tx.Users().Names().Get([]byte(name))
	userdata := tx.Users().Data().Get(userid)
	if userdata == nil {
		return nil
	}
	u := &User{}
	json.Unmarshal(userdata, u)
	return u
}

func (u *User) Save(t *bolt.Tx) {
	tx := &TX{t}
	if u.ID == nil {
		u.ID = NextKey(tx.Users().Data())
		tx.Users().Names().Put([]byte(u.Name), u.ID)
	}
	userdata, _ := json.Marshal(u)
	tx.Users().Data().Put(u.ID, userdata)
}

func (u *User) SetPassword(password string) {
	u.Password, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err == nil {
		return true
	}
	return false
}
