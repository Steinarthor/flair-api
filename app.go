package main

import (
	"database/sql"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type App struct {
	Router *mux.Router
	Db     *sql.DB
}

// Response type send to users.
type Response struct {
	Status  int
	Data    interface{}
	Message string
}

type Claims struct {
	Username string
	jwt.StandardClaims
}

// Initialize accepts database credentials and initializes the app.
func (a *App) Initialize(database string) {
	db, err := sql.Open(database, "./flair.db")
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()
	a.Db = db
	a.initializeRoutes()
}

// Run accepts an address and starts the application.
func (a *App) Run(addr string) {
	srv := &http.Server{
		Handler:      a.Router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Server is running and listening on port: 8081")
	log.Fatal(srv.ListenAndServe())
}

// Route handlers
func (a *App) AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var addUserRequest User

	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(w)
	err := dec.Decode(&addUserRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	err = addUserRequest.Add(addUserRequest.Name, addUserRequest.Username, addUserRequest.Email, a.Db)

	if err != nil {
		enc.Encode(Response{
			Status:  http.StatusConflict,
			Data:    "",
			Message: err.Error(),
		})

	} else {
		enc.Encode(Response{
			Status:  http.StatusCreated,
			Data:    addUserRequest,
			Message: "Success.",
		})
	}
}

func (a *App) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var getUserRequest User

	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(w)
	err := dec.Decode(&getUserRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	err, user := getUserRequest.Get(getUserRequest.Username, a.Db)

	if err != nil {
		enc.Encode(Response{
			Status:  http.StatusConflict,
			Data:    "",
			Message: err.Error(),
		})

	} else {
		enc.Encode(Response{
			Status:  http.StatusCreated,
			Data:    user,
			Message: "Success.",
		})
	}
}

func jwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" {
			http.Error(w, "Authorization token was not provided", http.StatusUnauthorized)
			return
		}

		extractedToken := strings.Split(token, "Bearer ")

		if len(extractedToken) == 2 {
			token = strings.TrimSpace(extractedToken[1])
		} else {
			http.Error(w, "Incorrect Format of Authorization Token", http.StatusBadRequest)
			return
		}

		parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SIGNING_KEY")), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Invalid Token Signature", http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !parsedToken.Valid {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/adduser", jwtMiddleware(a.AddUser)).Methods("POST")
	a.Router.HandleFunc("/getuser", jwtMiddleware(a.GetUser)).Methods("GET")
}
