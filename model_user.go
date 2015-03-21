package main

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"golang.org/x/crypto/bcrypt"
)

type Auth int

const (
	_                  = iota
	AuthMember    Auth = iota
	AuthModerator Auth = iota
	AuthAdmin     Auth = iota
)

func (a Auth) String() string {
	switch a {
	case AuthMember:
		return "Member"
	case AuthModerator:
		return "Moderator"
	case AuthAdmin:
		return "Admin"
	default:
		return ""
	}
}

type User struct {
	ID       []byte
	Name     string
	Password []byte
	Auths    []Auth
}

type SafeUser struct {
	Name  string
	Auths []Auth
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

func GetSafeUserByID(t *bolt.Tx, id []byte) *SafeUser {
	tx := TX{t}
	userdata := tx.Users().Data().Get(id)
	if userdata == nil {
		return nil
	}
	u := &SafeUser{}
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

func (u *User) HasAuth(auth Auth) bool {
	for _, a := range u.Auths {
		if a == auth {
			return true
		}
	}
	return false
}

func (u *User) GiveAuth(auth Auth) *User {
	if !u.HasAuth(auth) {
		u.Auths = append(u.Auths, auth)
	}
	return u
}
