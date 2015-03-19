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
	store    map[string]*Session
	mutex    sync.Mutex
	contexts map[*http.Request]*Session
}

var ()

func newMemoryStore() *MemoryStore {
	return &MemoryStore{
		store:    make(map[string]*Session),
		contexts: make(map[*http.Request]*Session),
	}
}

func (m *MemoryStore) Get(r *http.Request) *Session {
	if s, ok := m.contexts[r]; ok {
		return s
	}
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
	http.SetCookie(w, &http.Cookie{
		Name:    "GoWiki",
		Value:   s.ID,
		Path:    "/",
		Expires: time.Now().Add(time.Hour * 24 * 7),
	})
	return nil
}

func (m *MemoryStore) ContextClear(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			delete(m.contexts, req)
		}()
		h.ServeHTTP(rw, req)
	})
}
