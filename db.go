package auth

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type DB struct {
	Dir string
}

func read[T any](dir, t, id string) (*T, error) {
	path := filepath.Join(dir, t, id)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var v T
	err = json.NewDecoder(f).Decode(&v)
	if err != nil {
		panic(err)
	}
	return &v, nil
}

func save[T any](dir, t, id string, v *T) error {
	path := filepath.Join(dir, t, id)
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return os.WriteFile(path, b, os.ModePerm)
}

func (db *DB) User(id string) (*User, error) {
	return read[User](db.Dir, "users", id)
}

func (db *DB) SaveUser(u *User) error {
	return save(db.Dir, "users", u.ID, u)
}

func (db *DB) Phone(number string) (*Phone, error) {
	return read[Phone](db.Dir, "phones", number)
}

func (db *DB) SavePhone(p *Phone) error {
	return save(db.Dir, "phones", p.Number, p)
}

func (db *DB) Session(token string) (*Session, error) {
	return read[Session](db.Dir, "sessions", token)
}

func (db *DB) SaveSession(s *Session) error {
	return save(db.Dir, "sessions", s.Token, s)
}
