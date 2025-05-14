package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// Helper function to clear the database between tests
func clearDB() error {
	_, err := db.Exec("DELETE FROM users")
	return err
}

// Helper function to seed a test user
func seedTestUser(user User) error {
	_, err := db.Exec("INSERT INTO users (email, first_name, last_name) VALUES (?, ?, ?)",
		user.Email, user.FirstName, user.LastName)
	return err
}

func TestMain(m *testing.M) {
	// Set up test database
	var err error
	db, err = sql.Open("sqlite3", "./test.db")
	if err != nil {
		os.Exit(1)
	}
	defer os.Remove("./test.db")

	// Create tables
	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		email TEXT PRIMARY KEY,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL
	);`

	if _, err := db.Exec(createTable); err != nil {
		os.Exit(1)
	}

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestGetAllUsers(t *testing.T) {
	if err := clearDB(); err != nil {
		t.Fatalf("Error clearing database: %v", err)
	}

	// Seed a test user
	testUser := User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}
	if err := seedTestUser(testUser); err != nil {
		t.Fatalf("Error seeding test user: %v", err)
	}

	// Create a test request
	req := httptest.NewRequest("GET", "/example/get/users/all", nil)
	w := httptest.NewRecorder()

	// Call the handler
	getAllUsers(w, req)

	// Check the status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check the content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type %s, got %s", "application/json", contentType)
	}

	// Decode the response body
	var users []User
	if err := json.NewDecoder(w.Body).Decode(&users); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}
}

func TestGetUserByEmail(t *testing.T) {
	if err := clearDB(); err != nil {
		t.Fatalf("Error clearing database: %v", err)
	}

	// Seed a test user
	testUser := User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}
	if err := seedTestUser(testUser); err != nil {
		t.Fatalf("Error seeding test user: %v", err)
	}

	// Test case 1: User not found
	req := httptest.NewRequest("GET", "/example/get/user?email=notfound@example.com", nil)
	w := httptest.NewRecorder()

	getUserByEmail(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d for non-existent user, got %d", http.StatusNotFound, w.Code)
	}

	// Test case 2: Missing email parameter
	req = httptest.NewRequest("GET", "/example/get/user", nil)
	w = httptest.NewRecorder()

	getUserByEmail(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d for missing email, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateUser(t *testing.T) {
	if err := clearDB(); err != nil {
		t.Fatalf("Error clearing database: %v", err)
	}

	// Test case 1: Valid user creation
	user := User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/example/create/user", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	createUser(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Test case 2: Missing required fields
	invalidUser := User{
		FirstName: "John",
	}

	body, _ = json.Marshal(invalidUser)
	req = httptest.NewRequest("POST", "/example/create/user", bytes.NewBuffer(body))
	w = httptest.NewRecorder()

	createUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d for missing fields, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateUser(t *testing.T) {
	if err := clearDB(); err != nil {
		t.Fatalf("Error clearing database: %v", err)
	}

	// Seed a test user
	testUser := User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}
	if err := seedTestUser(testUser); err != nil {
		t.Fatalf("Error seeding test user: %v", err)
	}

	// Test case 1: Valid update
	updateData := User{
		Email:     "john.doe@example.com",
		FirstName: "Johnny",
	}

	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest("PUT", "/example/update/user", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	updateUser(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Test case 2: User not found
	notFoundUser := User{
		Email:     "notfound@example.com",
		FirstName: "NotFound",
	}

	body, _ = json.Marshal(notFoundUser)
	req = httptest.NewRequest("PUT", "/example/update/user", bytes.NewBuffer(body))
	w = httptest.NewRecorder()

	updateUser(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d for non-existent user, got %d", http.StatusNotFound, w.Code)
	}
}
