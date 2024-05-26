package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/chirps", createChirpHandler(db))
	mux.HandleFunc("GET /api/chirps", getChirpsHandler(db))
	mux.HandleFunc("/api/healthz", readinessHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}
}
