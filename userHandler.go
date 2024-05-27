package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)
func createUserHandler(db *DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost{
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return 
		}
		var req struct {
			Email string `json:"email"`
			Password string `json:"password"`
		}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Bad request")
			return
		}
		user, err := db.CreateUser(req.Email, req.Password)
		if err != nil {
			if err == ErrUserExists {
				respondWithError(w, http.StatusConflict, "User with this email already exists")
			} else {
				respondWithError(w, http.StatusInternalServerError, "Failed to create user")
			}
			return
		}
		response := struct{
			ID  int	`json:"id"`
			Email string `json:"email"`
		}{
			ID: user.ID,
			Email: user.Email,
		}

		respondWithJSON(w, http.StatusCreated, response);
	}
}
func loginHandler(db *DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost{
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		var req struct{
			Email string `json:"email"`
			Password string `json:"password"`
		}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Bad Request")
			return
		}

		user, exist := db.GetUserByEmail(req.Email)
		if !exist {
			respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		err = bcrypt.CompareHashAndPassword( []byte(user.HashedPassword) , []byte(req.Password))
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		response := struct{
			ID	int `json:"id"`
			Email	string `json:"email"` 
		}{
			ID: user.ID,
			Email: user.Email,
		}
		respondWithJSON(w, http.StatusOK, response)
	}
}
