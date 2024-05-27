package main

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}
type DBStructure struct {
	Users	map[int]User	`json:"users"`
	Chirps map[int]Chirp `json:"chirps"`
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
	data, err := os.ReadFile(db.path)
	if err != nil {
		if os.IsNotExist(err) {
			println("Database doesn't exist, creating empty database")
			initialData := DBStructure{
				Users: make(map[int]User),
				Chirps: make(map[int]Chirp),
			}
			return db.writeDB(initialData)
		}
		return err
	}
	// Check if the file content is empty or only contains whitespace
	if len(strings.TrimSpace(string(data))) == 0 {
		println("File is empty, initializing with empty chirps map");
		initialData := DBStructure{
			Users: make(map[int]User),
			Chirps: make(map[int]Chirp),
		}
		return db.writeDB(initialData)
	}
	// Check if the file contains an empty JSON object "{}"
	if string(data) == "{}" {
		println("File contains an empty JSON object, initializing with empty chirps map")
		initialData := DBStructure{
			Users: make(map[int]User),
			Chirps: make(map[int]Chirp),
		}
		return db.writeDB(initialData)
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
