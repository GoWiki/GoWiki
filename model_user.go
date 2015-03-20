package main

import (
	"encoding/json"
)

type User struct {
	ID []byte
	Name     string
	Password []byte
}

func GetUser(t *bolt.Tx, username string) *User, error {
	tx := &TX{t}
	user := &User{}
	userid := tx.Users().Names().Get([]byte(Name))
	userdata = tx.Users().Data().Get(userid)
	
	err := json.Unmarshal(userdata, user)
	if err != nil {
		return nil, err
	}
	return user, nil
	
}

func (u *User) SaveUser(t *bolt.Tx) {
	tx := &TX{t}
	b_userdata := tx.Users().Data()
	if u.ID == nil {
		u.ID = NextKey(b_userdata)
		tx.Users().Names().Set([]byte(u.Name), u.ID)
	}
	userdata, _ := json.Marshal(u)
	
	b_userdata.Set(u.ID, userdata)
}