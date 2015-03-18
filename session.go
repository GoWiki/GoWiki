package main

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type MemoryStore struct {
	Codecs  []securecookie.Codec
	Options *sessions.Options

	store map[string]map[interface{}]interface{}
	mutex sync.Mutex
}

var (
	ErrSessionNotFound error = errors.New("Session not found")
)

func newMemoryStore(keyPairs ...[]byte) *MemoryStore {
	return &MemoryStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &sessions.Options{
			Path:   "/",
			MaxAge: 86400 * 30,
		},
		store: make(map[string]map[interface{}]interface{}),
	}
}

func (m *MemoryStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(m, name)
}

func (m *MemoryStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(m, name)
	opts := *m.Options
	session.Options = &opts
	session.IsNew = true
	var err error
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, m.Codecs...)
		if err == nil {
			err = m.load(session)
			if err == nil {
				session.IsNew = false
			}
		}
	}
	return session, err
}

func (m *MemoryStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	if s.ID == "" {
		s.ID = string(securecookie.GenerateRandomKey(32))
	}
	if err := m.save(s); err != nil {
		return err
	}
	encoded, err := securecookie.EncodeMulti(s.Name(), s.ID,
		m.Codecs...)
	if err != nil {
		return err
	}
	http.SetCookie(w, sessions.NewCookie(s.Name(), encoded, s.Options))
	return nil
}

func (m *MemoryStore) load(session *sessions.Session) error {
	m.mutex.Lock()
	vals, ok := m.store[session.ID]
	m.mutex.Unlock()
	if ok {
		session.Values = vals
		return nil
	}
	return ErrSessionNotFound
}

func (m *MemoryStore) save(session *sessions.Session) error {
	m.mutex.Lock()
	m.store[session.ID] = session.Values
	m.mutex.Unlock()
	return nil
}
