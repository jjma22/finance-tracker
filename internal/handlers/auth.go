package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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
func (f *financeServer) VerifyUser(rw http.ResponseWriter, r *http.Request) {
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

	//user.Password, err = HashPassword(user.Password)

	userDetails, err := database.GetUser(user.Username)
	if err != nil {
		f.l.Error("Error getting username", "error", err)
		http.Error(rw, "User not found", http.StatusNotFound)
		return
	}

	if VerifyPassword(user.Password, userDetails.Password) != true {
		f.l.Error("Incorrect password", "error", err)
		http.Error(rw, "Username or password incorrect", http.StatusUnauthorized)
		return
	}

	f.l.Info("User login successful")

	// Generate token for authenticated user
	userToken, err := auth.CreateToken(user.Username)
	if err != nil {
		f.l.Error("Error generating token for user", "error", err)
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
	}

	// Write token into JSON
	resp, err := json.Marshal(&auth.UserToken{
		Token: userToken,
	})

	if err != nil {
		f.l.Error("Error marhsing usr token", "error", err)
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
	}

	// return token to user
	rw.Write(resp)

}

func (f *financeServer) CreateUser(rw http.ResponseWriter, r *http.Request) {
	// Get user details from request
	user, err := f.UserFromJSON(r)

	if err != nil {
		f.l.Error("Error decoding login request", "error", err)
		http.Error(rw, "Invalid request", http.StatusBadGateway)
		return
	}

	if user.Username == "" || user.Password == "" {
		f.l.Error("Invalid credentials sent to create user", "error", err)
		http.Error(rw, "Invalid username or password", http.StatusBadRequest)
		return
	}
	// Verify if user already exists

	exists, err := database.VerifyUserExists(user.Username)
	if err != nil {
		f.l.Error("Error checking if user exists", "error", err)
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
		return
	}

	if exists {
		f.l.Error("User already exists", "error", err)
		http.Error(rw, "User already exists", http.StatusBadRequest)
		return
	}
	// Generate hashed password
	hashedPw, err := HashPassword(user.Password)
	if err != nil {
		f.l.Error("Error hasing user password", "error", err)
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check orignal and hashed password match
	if VerifyPassword(user.Password, hashedPw) != true {
		f.l.Error("Hashed password does not match original request, not storing", "error", err)
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create new user in users table
	err = database.CreateUser(user.Username, hashedPw)
	if err != nil {
		f.l.Error("Error inserting new user into database", "error", err)
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
		return
	}

}

// Creates hash of password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (f *financeServer) MiddleWareValidateAuthentication(next http.Handler) http.Handler {
	// Annonymous function to validate expense before passing request onto next handler
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		authArr := strings.Split(authHeader, " ")

		// Errors if auth header splits into more than 2
		if len(authArr) != 2 {
			f.l.Error("Bearer token split into array of 3 during validation", "error", errors.New("Invalid authentication method"))
			http.Error(rw, "Invalid authentication method", http.StatusUnauthorized)
			return
		}
		// Checks Authorization method bearer is being used
		if authArr[0] != "Bearer" {
			f.l.Error("User did not use authentication method Bearer", "error", errors.New("Invalid authentication method"))
			http.Error(rw, "Unauthorized authentication method", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(authArr[1], func(token *jwt.Token) (any, error) {
			return []byte(auth.GetJwtKey()), nil
		})

		if err != nil {
			f.l.Error("Error parsing token", "error", err)
			http.Error(rw, "Error authorizing credentials", http.StatusInternalServerError)
			return
		}

		switch {
		case token.Valid:
			next.ServeHTTP(rw, r)
		case errors.Is(err, jwt.ErrTokenMalformed):
			f.l.Error("User using malformed token", "error", errors.New("Invalid token"))
			http.Error(rw, "Error authorizing credentials", http.StatusUnauthorized)
			return
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			// Invalid signature
			f.l.Error("User using token with invalid signature", "error", errors.New("Invalid token"))
			http.Error(rw, "Error authorizing credentials", http.StatusUnauthorized)
			return
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			// Token is either expired or not active yet
			f.l.Error("User token expired", "error", errors.New("Invalid token"))
			http.Error(rw, "Error authorizing credentials", http.StatusUnauthorized)
			return
		default:
			f.l.Error("Error handling user token", "error", errors.New("Invalid token"))
			http.Error(rw, "Error authorizing credentials", http.StatusUnauthorized)
			return
		}
	})

}
