package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jjma22/finance-tracker/internal/auth"
)

func (f *financeServer) LoginUser(rw http.ResponseWriter, r *http.Request) {
	var user auth.User

	// read json in user
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		f.l.Error("Error decoding login request", "error", err)
		http.Error(rw, "Invalid request", http.StatusBadGateway)
	}

	// Checks for empty username
	if user.Username == "" {
		f.l.Error("Invalid credentials sent to login", "error", errors.New("Invalid or empty username"))
		http.Error(rw, "Invalid username", http.StatusBadGateway)

	}

	// Checks for empty password
	if user.Username == "" {
		f.l.Error("Invalid credentials sent to login", "error", errors.New("Invalid or empty password"))
		http.Error(rw, "Invalid password", http.StatusBadGateway)
	}

}
