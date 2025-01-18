package auth

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type DB struct {
	Dir string
}

func (db *DB) User(id string) (*User, error) {
	path := filepath.Join(db.Dir, "users", id)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	u := User{}
	err = json.NewDecoder(f).Decode(&u)
	if err != nil {
		panic(err)
	}
	return &u, nil
}

func (db *DB) Phone(number string) (*Phone, error) {
	path := filepath.Join(db.Dir, "phones", number)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	p := Phone{}
	err = json.NewDecoder(f).Decode(&p)
	if err != nil {
		panic(err)
	}
	return &p, nil
}
