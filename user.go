package main
type User struct{
	ID int	`json:"id"`
	Email string	`json:"email"`
}
func (db* DB) CreateUser(email string)(User, error){
	db.mux.Lock()
	defer db.mux.Unlock()
	
	DBStructure, err := db.loadDB()
	if err != nil {
		return User{},err
	}
	newID := len(DBStructure.Users) + 1;
	user := User{
		ID: newID,
		Email: email,
	}
	DBStructure.Users[newID] = user;
	if err = db.writeDB(DBStructure); err != nil{
		return User{}, err;
	}
	return user, nil;
}
func (db *DB) GetUserByID(id int) (User, bool) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, false
	}

	user, exists := dbStructure.Users[id]
	return user, exists
}