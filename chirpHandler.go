package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

// Handler for creating a new chirp
func createChirpHandler(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		var req struct {
			Body string `json:"body"`
			UID int	`json:"uid"`
		}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Bad request")
			return
		}

		if len(req.Body) > 140 {
			respondWithError(w, http.StatusBadRequest, "Chirp is too long")
			return
		}

		cleanedBody := replaceProfaneWords(req.Body)
		chirp, err := db.CreateChirp(cleanedBody, req.UID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create chirp")
			return
		}

		respondWithJSON(w, http.StatusCreated, chirp)
	}
}

// Handler for getting all chirps
func getChirpsHandler(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		chirps, err := db.GetChirps()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to get chirps")
			return
		}

		respondWithJSON(w, http.StatusOK, chirps)
	}
}
func getChirpByIDHandler(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		// Extract the path parameters
		idStr := chi.URLParam(r, "id")

		// Convert the ID to an integer
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
			return
		}

		// Get the chirp by ID
		chirp, exists := db.GetChirpByID(id)
		if !exists {
			respondWithError(w, http.StatusNotFound, "Chirp not found")
			return
		}

		// Respond with the chirp
		respondWithJSON(w, http.StatusOK, chirp)
	}
}