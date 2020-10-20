package main

import (
	"database/sql"
	"errors"
	"log"
)

// User model
type User struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Signup accepts username, password and an email and registers an creates a new account.
func (u *User) Add(name, username, email string, db *sql.DB) error {
	addUser, err := db.Prepare("INSERT INTO user(name, username, email) VALUES(?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = addUser.Exec(name, username, email)

	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (u *User) Get(username string, db *sql.DB) (error, User) {
	var user User

	getUser, err := db.Prepare("SELECT name, username, email FROM user WHERE username = ?")
	if err != nil {
		log.Fatal(err)
	}

	err = getUser.QueryRow(username).Scan(&user.Name, &user.Username, &user.Email)

	if err != nil {
		return errors.New("no user found"), User{}
	}
	return nil, user
}
