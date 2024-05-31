package main

import (
	"errors"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)
type User struct{
	ID int	`json:"id"`
	HashedPassword string 
	Email string	`json:"email"`
}
func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	if !regexp.MustCompile(`[\W]`).MatchString(password) {
		return errors.New("password must contain at least one special character")
	}

	return nil
}
func (db* DB) CreateUser(email string, password string)(User,int, error){
	db.mux.Lock()
	defer db.mux.Unlock()
	if err:= validatePassword(password); err != nil {
		return User{}, http.StatusBadRequest,err
	}
	DBStructure, err := db.loadDB()
	if err != nil {
		return User{},http.StatusInternalServerError,err
	}
	for _, user := range DBStructure.Users{
		if user.Email == email {
			return User{}, http.StatusConflict,ErrUserExists
		}
	}
	newID := len(DBStructure.Users) + 1;
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost);
	if err != nil {
		return User{}, http.StatusInternalServerError,err
	}
	user := User{
		ID: newID,
		HashedPassword: string(hashedPassword),
		Email: email,
	}
	DBStructure.Users[newID] = user;
	if err = db.writeDB(DBStructure); err != nil{
		return User{}, http.StatusInternalServerError,err;
	}
	return user, http.StatusCreated, nil;
}
func (db *DB) UpdateUser(user User, email, password string)(int, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	
	DBStructure, err := db.loadDB()
	if err != nil {
		return http.StatusInternalServerError,err
	}
	
	if email != "" {
		user.Email = email
	}
	if password != "" {
		if err := validatePassword(password); err != nil {
			return http.StatusBadRequest, err
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return http.StatusInternalServerError, errors.New("failed to hash password")
		}
		user.HashedPassword = string(hashedPassword)
	}

	DBStructure.Users[user.ID] = user
	if err = db.writeDB(DBStructure); err != nil{
		return http.StatusInternalServerError,err;
	}
	return http.StatusAccepted, nil;
}
func (db *DB) GetUserByID(userID int)(User, bool){
	db.mux.Lock()
	defer db.mux.Unlock()

	DBStructure, err := db.loadDB()
	
	if err != nil {
		return User{}, false
	}
	for _, user := range DBStructure.Users{
		if user.ID == userID {
			return user, true
		}
	}
	return User{}, false
}
func (db *DB) GetUserByEmail(email string) (User, bool) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure, err := db.loadDB()
	
	if err != nil {
		return User{}, false
	}
	for _, user := range dbStructure.Users{
		if user.Email == email {
			return user, true
		}
	}
	return User{}, false
}
