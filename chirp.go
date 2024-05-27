package main

import "sort"
type Chirp struct {
	ID int `json:"id"`
	Body string `json:"body"`
	UID int `json:"uid"`
}
func (db *DB) CreateChirp(body string,uid int) (Chirp, error) {
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
		UID: uid,
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
func (db *DB) GetChirpByID(id int) (Chirp, bool) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, false
	}

	chirp, exists := dbStructure.Chirps[id]
	return chirp, exists
}