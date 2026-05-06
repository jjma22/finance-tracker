package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jjma22/finance-tracker/internal/auth"
	"github.com/jjma22/finance-tracker/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (f *financeServer) UserFromJSON(r *http.Request) (*auth.User, error) {
	var u auth.User

	// read json in user
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		return nil, err
	}

	return &u, nil

}
func (f *financeServer) LoginUser(rw http.ResponseWriter, r *http.Request) {
	// Get user details from request
	user, err := f.UserFromJSON(r)

	if err != nil {
		f.l.Error("Error decoding login request", "error", err)
		http.Error(rw, "Invalid request", http.StatusBadGateway)
		return
	}

	// Checks for empty username
	if user.Username == "" {
		f.l.Error("Invalid credentials sent to login", "error", errors.New("Invalid or empty username"))
		http.Error(rw, "Invalid username", http.StatusBadGateway)
		return

	}

	// Checks for empty password
	if user.Password == "" {
		f.l.Error("Invalid credentials sent to login", "error", errors.New("Invalid or empty password"))
		http.Error(rw, "Invalid password", http.StatusBadGateway)
		return
	}

	user.Password, err = HashPassword(user.Password)

	userDetails, err := database.GetUser(user.Username)
	if err != nil {
		f.l.Error("Error getting username", "error", err)
		http.Error(rw, "User not found", http.StatusNotFound)
		return
	}

	fmt.Println(userDetails)

}

func (f *financeServer) CreateUser(rw http.ResponseWriter, r *http.Request) {
	// Get user details from request
	user, err := f.UserFromJSON(r)

	if err != nil {
		f.l.Error("Error decoding login request", "error", err)
		http.Error(rw, "Invalid request", http.StatusBadGateway)
		return
	}

	hashedPw, err := HashPassword(user.Password)
	if err != nil {
		f.l.Error("Error hasing user password", "error", err)
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
		return
	}

	if VerifyPassword(user.Password, hashedPw) != true {
		f.l.Error("Hashed password does not match original request, not storing", "error", err)
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
		return
	}

}
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
