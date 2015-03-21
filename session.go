package main

import (
	"net/http"
	"sync"
	"time"
)

type Session struct {
	ID                string
	User              *User
	LastAccess        *time.Time
	PostLoginRedirect string
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
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if s, ok := m.contexts[r]; ok {
		return s
	}
	c, err := r.Cookie("GoWiki")
	if err != nil {
		s := &Session{}
		s.ID = GetRandomID()
		m.store[s.ID] = s
		return s
	}
	s, ok := m.store[c.Value]
	if !ok {
		s := &Session{}
		s.ID = GetRandomID()
		m.store[s.ID] = s
		return s
	}
	return s
}

func (m *MemoryStore) Destroy(r *http.Request, w http.ResponseWriter, s *Session) {
	delete(m.store, s.ID)
	delete(m.contexts, r)
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
