package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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
		user, code, err := db.CreateUser(req.Email, req.Password)
		if err != nil {
			respondWithError(w, code, err.Error())
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
func loginHandler(apiConfig *ApiConfig) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost{
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		var req struct{
			Email string `json:"email"`
			Password string `json:"password"`
			ExpiresInSeconds int `json:"expires_in_seconds"`
		}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Bad Request")
			return
		}

		user, exist := apiConfig.db.GetUserByEmail(req.Email)
		if !exist {
			respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		err = bcrypt.CompareHashAndPassword( []byte(user.HashedPassword) , []byte(req.Password))
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		
		tokenString, err := generateJWT(user.ID, apiConfig.rsaPrivate, int64(req.ExpiresInSeconds))
		if err != nil{
			respondWithError(w, http.StatusInternalServerError, "Failed to generate Token")
			return
		}
		response := struct{
			ID	int `json:"id"`
			Email	string `json:"email"` 
			Token 	string `json:"token"`
		}{
			ID: user.ID,
			Email: user.Email,
			Token: tokenString,
		}
		respondWithJSON(w, http.StatusOK, response)
	}
}
func updateUserHandler(apiConfig *ApiConfig) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			respondWithError(w, http.StatusMethodNotAllowed, "Wrong method")
			return
		}

		headerToken := r.Header.Get("Authorization")
		tokenString := strings.TrimSpace(strings.TrimPrefix(headerToken, "Bearer "))
        if tokenString == "" {
            respondWithError(w, http.StatusUnauthorized, "Invalid token format")
            return
        }
		token, err := jwt.ParseWithClaims(headerToken, &jwt.RegisteredClaims{} ,func(t *jwt.Token) (interface{}, error) {return apiConfig.rsaPublic, nil})
		println(token.Valid)
		if err != nil{
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid{
			userID := claims.Subject
			var req struct {
				Email	string `json:"email"`
				Password	string `json:"password"`
			}
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&req); err!= nil {
				respondWithError(w,http.StatusBadRequest, "Bad request")
				return
			}
			userIDInt, err := strconv.Atoi(userID)
			if err != nil {
				respondWithError(w,http.StatusBadRequest, "Invalid user ID")
				return
			}
			user, exist := apiConfig.db.GetUserByID(userIDInt)
			if !exist {
				respondWithError(w, http.StatusNotFound, "User not found")
				return
			}

			if code, err:= apiConfig.db.UpdateUser(user, req.Email, req.Password); err != nil{
				respondWithError(w, code, err.Error())
				return
			}
			respondWithJSON(w, http.StatusOK, "Credentials updated!")
		}else{
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		}
	}

}