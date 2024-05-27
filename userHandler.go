package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)
func createUserHandler(db *DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost{
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return 
		}
		var req struct {
			Email string `json:"email"`
		}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Bad request")
			return
		}
		user, err := db.CreateUser(req.Email)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create User")
			return
		}
		respondWithJSON(w, http.StatusCreated, user);
	}
}
func GetUserHandler(db *DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet{
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		strId := chi.URLParam(r, "uid")
		id,err := strconv.Atoi(strId);
		if err != nil {
			respondWithError(w,http.StatusBadRequest, "Wrong format of ID")
			return
		}
		user, exist := db.GetUserByID(id)
		if !exist{
			respondWithError(w,http.StatusNotFound, "user ID not found")
			return
		}
		respondWithJSON(w,http.StatusOK, user.Email)
	}
}