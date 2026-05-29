package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jjma22/finance-tracker/internal/auth"
	env_config "github.com/jjma22/finance-tracker/internal/config"
	"github.com/jjma22/finance-tracker/internal/database"
	"github.com/jjma22/finance-tracker/internal/handlers"
)

var TestUser = auth.User{
	Username: "test_user",
	Password: "test_abc!£$%",
}

// Want to create shared function for all handlers
func initDBTestAuth() {
	Config := *env_config.LoadConfig("../../.env")

	db_connection := Config.Database

	// Declare logger
	l := slog.Default()

	// Setup database connection
	database.InitDb(l, &db_connection)
}

// Tests user can be created successfully
func TestCreateUser(t *testing.T) {
	t.Run("creates user", func(t *testing.T) {

		en, err := json.Marshal(TestUser)

		if err != nil {
			t.Fatalf("Unable to parse user %v , '%v'", TestUser, err)
		}
		request, _ := http.NewRequest(http.MethodPost, "/create/user", bytes.NewReader(en))
		response := httptest.NewRecorder()

		l := slog.Default()
		initDBTestAuth()

		fh := handlers.FinanceNewServer(l)

		fh.CreateUser(response, request)

		want := 200
		got := response.Code

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}

	})
}

// Tests creating a user with username that already exists
func TestCreateUsernameAlreadyExists(t *testing.T) {
	t.Run("creates user", func(t *testing.T) {

		en, err := json.Marshal(TestUser)

		if err != nil {
			t.Fatalf("Unable to parse user %v , '%v'", TestUser, err)
		}
		request, _ := http.NewRequest(http.MethodPost, "/create/user", bytes.NewReader(en))
		response := httptest.NewRecorder()

		l := slog.Default()
		initDBTestAuth()

		fh := handlers.FinanceNewServer(l)

		fh.CreateUser(response, request)

		want := 400
		got := response.Code

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}

		expectedError := "User already exists"
		returnedError := strings.TrimSpace(response.Body.String())

		if returnedError != expectedError {
			t.Errorf("got '%s', want '%s'", returnedError, expectedError)
		}
	})
}

// Tests creating a user fails is username is not specified
func TestCreateUserNoPasswordSupplied(t *testing.T) {
	t.Run("creates user", func(t *testing.T) {

		u := &auth.User{
			Username: "test_user",
		}

		en, err := json.Marshal(u)

		if err != nil {
			t.Fatalf("Unable to parse user %v , '%v'", u, err)
		}
		request, _ := http.NewRequest(http.MethodPost, "/create/user", bytes.NewReader(en))
		response := httptest.NewRecorder()

		l := slog.Default()
		initDBTestAuth()

		fh := handlers.FinanceNewServer(l)

		fh.CreateUser(response, request)

		want := 400
		got := response.Code

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}

	})
}

// Hashed password returns different value to original password, but can be verified with the same password
func TestHashedPassword(t *testing.T) {
	t.Run("Test password is hashed", func(t *testing.T) {

		hashedPw, err := handlers.HashPassword(TestUser.Password)
		if err != nil {
			t.Fatalf("Error hashing password %v , '%v'", TestUser.Password, err)
		}

		if hashedPw == TestUser.Password {
			t.Errorf("Hashed password should not be the same as original password")
		}

		if handlers.VerifyPassword(TestUser.Password, hashedPw) != true {
			t.Errorf("Original password should match hashed password")
		}
	})
}

func TestVerifyUser(t *testing.T) {
	t.Run("Test user is successfully authenticated", func(t *testing.T) {

		en, err := json.Marshal(TestUser)

		if err != nil {
			t.Fatalf("Unable to parse user %v , '%v'", TestUser, err)
		}
		request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(en))
		response := httptest.NewRecorder()

		l := slog.Default()
		initDBTestAuth()

		fh := handlers.FinanceNewServer(l)

		fh.VerifyUser(response, request)

		want := 200
		got := response.Code

		if got != want {
			t.Errorf("got %d, want %d", got, want)

		}

		var token auth.UserToken

		err = json.Unmarshal(response.Body.Bytes(), &token)

		if err != nil {
			t.Errorf("Unable to parse response from server %v into user token, '%v'", response.Body, err)
		}

	})
}

func TestVerifyUserIncorrectPw(t *testing.T) {
	t.Run("Test user authentication fails with incorrect password", func(t *testing.T) {

		testUserPw := &auth.User{
			Username: TestUser.Username,
			Password: "12345",
		}
		fmt.Println(testUserPw)
		en, err := json.Marshal(testUserPw)

		if err != nil {
			t.Fatalf("Unable to parse user %v , '%v'", testUserPw, err)
		}
		request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(en))
		response := httptest.NewRecorder()

		l := slog.Default()
		initDBTestAuth()

		fh := handlers.FinanceNewServer(l)

		fh.VerifyUser(response, request)

		want := 401
		got := response.Code

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}

	})
}

func TestVerifyUserNoPw(t *testing.T) {
	t.Run("Test user authentication fails with no password", func(t *testing.T) {

		testUserPw := &auth.User{
			Username: TestUser.Username,
		}
		fmt.Println(testUserPw)
		en, err := json.Marshal(testUserPw)

		if err != nil {
			t.Fatalf("Unable to parse user %v , '%v'", testUserPw, err)
		}
		request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(en))
		response := httptest.NewRecorder()

		l := slog.Default()
		initDBTestAuth()

		fh := handlers.FinanceNewServer(l)

		fh.VerifyUser(response, request)

		want := 502
		got := response.Code

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}

	})
}

func TestVerifyUserNoUsername(t *testing.T) {
	t.Run("Test user authentication fails with no username", func(t *testing.T) {

		testUserPw := &auth.User{
			Password: "password",
		}
		fmt.Println(testUserPw)
		en, err := json.Marshal(testUserPw)

		if err != nil {
			t.Fatalf("Unable to parse user %v , '%v'", testUserPw, err)
		}
		request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(en))
		response := httptest.NewRecorder()

		l := slog.Default()
		initDBTestAuth()

		fh := handlers.FinanceNewServer(l)

		fh.VerifyUser(response, request)

		want := 502
		got := response.Code

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}

	})
}
