package main

import (
	"database/sql"
	"errors"
	"log"
)

// User model
type User struct {
	ID      int64 `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

// Signup accepts, password and an email and registers an creates a new account.
func (u *User) Add(name, email string, db *sql.DB) error {
	addUser, err := db.Prepare("INSERT INTO user(name, email) VALUES(?,?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = addUser.Exec(name, email)

	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (u *User) Get(email string, db *sql.DB) (error, User) {
	var user User

	getUser, err := db.Prepare("SELECT id, name, email FROM user WHERE email = ?")
	if err != nil {
		log.Fatal(err)
	}

	err = getUser.QueryRow(email).Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		return errors.New("no user found"), User{}
	}
	return nil, user
}
