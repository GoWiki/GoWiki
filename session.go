package main

import (
	"net/http"
	"sync"
	"time"
)

type Session struct {
	ID         string
	User       *User
	LastAccess *time.Time
}

type MemoryStore struct {
	store map[string]Session
	mutex sync.Mutex
}

var ()

func newMemoryStore() *MemoryStore {
	return &MemoryStore{
		store: make(map[string]*Session),
	}
}

func (m *MemoryStore) Get(r *http.Request) (*Session, error) {
	c, err := r.Cookie("GoWiki")
	if err != nil {
		return &Session{}
	}
	s, ok := m.store[c.Value]
	if !ok {
		s = &Session{}
	}
	return s
}

func (m *MemoryStore) Save(r *http.Request, w http.ResponseWriter, s *Session) error {

}
