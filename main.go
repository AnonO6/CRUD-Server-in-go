package main

import (
	"crypto/rsa"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

type ApiConfig struct{
	db	*DB
	jwtSecret	string
	rsaPrivate	*rsa.PrivateKey
	rsaPublic 	*rsa.PublicKey
}
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

// For pinging the server
func readinessHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Add("content-type","text/plain; charset=utf-8");
	w.WriteHeader(http.StatusOK);
	w.Write([]byte(http.StatusText(http.StatusOK)));
}

func main() {
	// We dont need jwt secret right now since we are using RSA
	err := godotenv.Load();
	if err != nil{
		log.Fatalf("Failed to load environment variable: %s",err)
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET is not set")
	}
	db, err := NewDB("database.json")
	if err != nil {
		log.Fatalf("Failed to initialize database: %s", err)
		return
	}
	privateKey, publicKey, err := loadRSAKeys()
	if err != nil {
		log.Fatalf("Failed to load RSA keys: %v", err)
	}
	
	apiConfig := &ApiConfig{
		db: db,
		jwtSecret: jwtSecret,
		rsaPrivate: privateKey,
		rsaPublic: publicKey,
	}

	r := chi.NewRouter()

	// Some middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Defining the routes
	r.Get("/api/healthz",readinessHandler)
	// For chirps
	r.Post("/api/chirps", createChirpHandler(db))
	r.Get("/api/chirps", getChirpsHandler(db))
	r.Get("/api/chirps/{id}", getChirpByIDHandler(db))
	// For Users
	r.Post("/api/users", createUserHandler(db))
	r.Put("/api/users", updateUserHandler(apiConfig))
	r.Post("/api/login", loginHandler(apiConfig))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("could not start server: %s\n", err)
		return
	}
}
