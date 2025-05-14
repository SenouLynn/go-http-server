package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		return err
	}

	// Create users table if it doesn't exist
	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		email TEXT PRIMARY KEY,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL
	);`

	_, err = db.Exec(createTable)
	return err
}

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Go HTTP Server!")
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT first_name, last_name, email FROM users")
	if err != nil {
		http.Error(w, "Error fetching users from database", http.StatusInternalServerError)
		log.Printf("Error fetching users: %v", err)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.FirstName, &user.LastName, &user.Email)
		if err != nil {
			http.Error(w, "Error reading user data", http.StatusInternalServerError)
			log.Printf("Error scanning user row: %v", err)
			return
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error reading user data", http.StatusInternalServerError)
		log.Printf("Error after scanning users: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func getUserByEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email parameter is required", http.StatusBadRequest)
		return
	}

	var user User
	err := db.QueryRow("SELECT first_name, last_name, email FROM users WHERE email = ?", email).
		Scan(&user.FirstName, &user.LastName, &user.Email)
	
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		log.Printf("Error fetching user: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if user.FirstName == "" || user.LastName == "" || user.Email == "" {
		http.Error(w, "FirstName, LastName, and Email are required", http.StatusBadRequest)
		return
	}

	// Insert user into database
	_, err = db.Exec("INSERT INTO users (email, first_name, last_name) VALUES (?, ?, ?)",
		user.Email, user.FirstName, user.LastName)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		log.Printf("Error creating user: %v", err)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	// Validate that email is provided
	if user.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Validate that at least first name or last name is provided
	if user.FirstName == "" && user.LastName == "" {
		http.Error(w, "Either firstName or lastName must be provided", http.StatusBadRequest)
		return
	}

	// Check if user exists
	var existingUser User
	err = db.QueryRow("SELECT first_name, last_name, email FROM users WHERE email = ?", user.Email).
		Scan(&existingUser.FirstName, &existingUser.LastName, &existingUser.Email)
	
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		log.Printf("Error fetching user: %v", err)
		return
	}

	// Update user in database
	if user.FirstName == "" {
		user.FirstName = existingUser.FirstName
	}
	if user.LastName == "" {
		user.LastName = existingUser.LastName
	}

	_, err = db.Exec("UPDATE users SET first_name = ?, last_name = ? WHERE email = ?",
		user.FirstName, user.LastName, user.Email)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		log.Printf("Error updating user: %v", err)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func main() {
	if err := initDB(); err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Close()

	http.HandleFunc("/", homePage)
	http.HandleFunc("/example/get/users/all", getAllUsers)
	http.HandleFunc("/example/get/user", getUserByEmail)
	http.HandleFunc("/example/create/user", createUser)
	http.HandleFunc("/example/update/user", updateUser)
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
