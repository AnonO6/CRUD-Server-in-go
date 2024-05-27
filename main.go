package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Define the list of profane words
var profaneWords = []string{"fuck", "shit", "chutiya"}

// Function to replace profane words
func replaceProfaneWords(text string) string {
	words := strings.Fields(text)
	for i, word := range words {
		cleanWord := strings.Trim(word, ".,!?;")
		for _, profane := range profaneWords {
			if strings.EqualFold(strings.ToLower(cleanWord), strings.ToLower(profane)) {
				words[i] = "****"
				break
			}
		}
	}
	return strings.Join(words, " ")
}

// Helper function to respond with error
func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// Helper function to respond with JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// Handler for creating a new chirp
func createChirpHandler(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		var req struct {
			Body string `json:"body"`
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
		chirp, err := db.CreateChirp(cleanedBody)
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
func readinessHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Add("content-type","text/plain; charset=utf-8");
	w.WriteHeader(http.StatusOK);
	w.Write([]byte(http.StatusText(http.StatusOK)));
}

func main() {
	db, err := NewDB("database.json")
	if err != nil {
		log.Fatalf("Failed to initialize database: %s", err)
	}
	r := chi.NewRouter()

	// Use some middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Define the routes
	r.Post("/api/chirps", createChirpHandler(db))
	r.Get("/api/chirps", getChirpsHandler(db))
	r.Get("/api/chirps/{id}", getChirpByIDHandler(db))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}
}
