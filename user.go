package main

import (
	"golang.org/x/crypto/bcrypt"
)
type User struct{
	ID int	`json:"id"`
	HashedPassword string 
	Email string	`json:"email"`
}
func (db* DB) CreateUser(email string, password string)(User, error){
	db.mux.Lock()
	defer db.mux.Unlock()
	
	DBStructure, err := db.loadDB()
	if err != nil {
		return User{},err
	}
	for _, user := range DBStructure.Users{
		if user.Email == email {
			return User{}, ErrUserExists
		}
	}
	newID := len(DBStructure.Users) + 1;
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost);
	if err != nil {
		return User{}, err
	}
	user := User{
		ID: newID,
		HashedPassword: string(hashedPassword),
		Email: email,
	}
	DBStructure.Users[newID] = user;
	if err = db.writeDB(DBStructure); err != nil{
		return User{}, err;
	}
	return user, nil;
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
