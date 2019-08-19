package session

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var GlobalSessionStore SessionStore

type SessionStore interface {
	Find(string) (*Session, error)
	Save(*Session) error
	Delete(*Session) error
}

type FileSessionStore struct {
	filename string
	Session  map[string]Session
}

func NewFileSessionStore(name string) (*FileSessionStore, error) {
	store := &FileSessionStore{
		Session:  map[string]Session{},
		filename: name,
	}

	content, err := ioutil.ReadFile(name)

	if err != nil {
		//If it is a matter of the file not existing, that's ok
		if os.IsNotExist(err) {
			return store, nil
		}
		return nil, err
	}

	err = json.Unmarshal(content, store)
	if err != nil {
		return nil, err
	}
	return store, err
}

func (s *FileSessionStore) Find(id string) (*Session, error) {
	session, exists := s.Session[id]
	if !exists {
		return nil, nil
	}
	return &session, nil
}

func (store *FileSessionStore) Save(session *Session) error {
	store.Session[session.ID] = *session
	contents, err := json.MarshalIndent(store, "", "")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(store.filename, contents, 0660)
}

func (store *FileSessionStore) Delete(session *Session) error {
	delete(store.Session, session.ID)
	content, err := json.MarshalIndent(store, "", "")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(store.filename, content, 0660)
}
