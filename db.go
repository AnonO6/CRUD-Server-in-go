package main

import (
	"encoding/json"
	"os"
	"sort"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	err := db.ensureDB()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if err != nil {
		if os.IsNotExist(err) {
			initialData := DBStructure{Chirps: make(map[int]Chirp)}
			return db.writeDB(initialData)
		}
		return err
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	var dbStructure DBStructure
	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		return DBStructure{}, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	data, err := json.MarshalIndent(dbStructure, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(db.path, data, 0644)
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newID := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:   newID,
		Body: body,
	}

	dbStructure.Chirps[newID] = chirp
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	var chirps []Chirp
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	// Sort chirps by ID
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	return chirps, nil
}
