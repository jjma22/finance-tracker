package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jjma22/finance-tracker/internal/auth"
	env_config "github.com/jjma22/finance-tracker/internal/config"
	"github.com/jjma22/finance-tracker/internal/data"
	"github.com/jjma22/finance-tracker/internal/database"
	"github.com/jjma22/finance-tracker/internal/handlers"
)

var TestUser = auth.User{
	Username: "test_user",
	Password: "test_abc!£$%",
}

var TestUserToken auth.UserToken

var TestUserClaims jwt.MapClaims

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

		err = json.Unmarshal(response.Body.Bytes(), &User2Token)

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

var User2 auth.User

var User2Token auth.UserToken

var User2Claims jwt.MapClaims

func TestCannotGetOtherUserBudget(t *testing.T) {
	t.Run("Can create second user and add budget", func(t *testing.T) {
		User2 = auth.User{
			Username: "test2",
			Password: "password1234",
		}

		u2, err := json.Marshal(User2)

		if err != nil {
			t.Fatalf("Unable to parse user %v , '%v'", User2, err)
		}

		request, _ := http.NewRequest(http.MethodPost, "/create/user", bytes.NewReader(u2))
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

		request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewReader(u2))

		fh.VerifyUser(response, request)

		err = json.Unmarshal(response.Body.Bytes(), &User2Token)

		if err != nil {
			t.Errorf("Unable to parse response from server %v into user token, '%v'", response.Body, err)
		}

		user2Budget := &data.Budget{
			Budget: 2000,
		}

		request, _ = http.NewRequest(http.MethodPost, "/monthlybudget", nil)
		ctx := context.WithValue(request.Context(), handlers.Budget{}, user2Budget)

		user2token, err := jwt.Parse(User2Token.Token, func(token *jwt.Token) (any, error) {
			return []byte(auth.GetJwtKey()), nil
		})

		if err != nil {
			t.Fatalf("Error parsing token, error: %v", err)
		}

		User2Claims = user2token.Claims.(jwt.MapClaims)

		//Set userid in context for handlers
		ctx = context.WithValue(ctx, handlers.UserKey{}, User2Claims["username"])

		request = request.WithContext(ctx)
		fh.SetBudget(response, request)

		want = 200
		got = response.Code

		if got != want {
			t.Errorf("got %d, want %d when adding budget for User2", got, want)
		}
	})

	t.Run("Second user can get new budget", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/monthlybudget/1", nil)
		response := httptest.NewRecorder()

		l := slog.Default()
		initDBTestAuth()
		// Manualy inject path
		request.SetPathValue("id", "1")
		fh := handlers.FinanceNewServer(l)
		ctx := context.WithValue(request.Context(), handlers.UserKey{}, User2Claims["username"])

		request = request.WithContext(ctx)

		fh.GetBudget(response, request)

		var got data.Budget
		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse response from server %d into budget, '%v'", response.Body, err)
		}
		want := 2000

		if got.Budget != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("Test user cannot get user2 budget", func(t *testing.T) {
		testUser, err := json.Marshal(TestUser)

		if err != nil {
			t.Fatalf("Unable to parse user %v , '%v'", User2, err)
		}

		request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(testUser))
		response := httptest.NewRecorder()

		l := slog.Default()
		initDBTestAuth()

		fh := handlers.FinanceNewServer(l)
		fh.VerifyUser(response, request)

		err = json.Unmarshal(response.Body.Bytes(), &TestUserToken)

		if err != nil {
			t.Errorf("Unable to parse response from server %v into user token, '%v'", response.Body, err)
		}

		request, _ = http.NewRequest(http.MethodGet, "/monthlybudget/1", nil)
		response = httptest.NewRecorder()

		// Manualy inject path
		request.SetPathValue("id", "1")

		testUserToken, err := jwt.Parse(TestUserToken.Token, func(token *jwt.Token) (any, error) {
			return []byte(auth.GetJwtKey()), nil
		})

		if err != nil {
			t.Fatalf("Error parsing token, error: %v", err)
		}

		TestUserClaims = testUserToken.Claims.(jwt.MapClaims)
		fmt.Println(TestUserClaims)

		ctx := context.WithValue(request.Context(), handlers.UserKey{}, TestUserClaims["username"])

		request = request.WithContext(ctx)

		fh.GetBudget(response, request)

		want := 500

		if response.Code != want {
			t.Errorf("got %d, want %d", response.Code, want)
		}
	})

	t.Run("User can update budget", func(t *testing.T) {
		b, err := json.Marshal(&data.Budget{
			Budget: 1500,
		})

		if err != nil {
			t.Fatalf("Unable to parse budget from client %d , '%v'", b, err)
		}
		request, _ := http.NewRequest(http.MethodPut, "/monthlybudget/1", bytes.NewReader(b))
		response := httptest.NewRecorder()

		l := slog.Default()
		initDBTestAuth()
		// Manualy inject path
		request.SetPathValue("id", "1")
		fh := handlers.FinanceNewServer(l)

		uuid := User2Claims["username"]
		ctx := context.WithValue(request.Context(), handlers.UserKey{}, uuid)

		request = request.WithContext(ctx)
		fh.UpdateBudget(response, request)

		want := 200
		got := response.Code
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}

		// Check id 1 budget has updated
		request, _ = http.NewRequest(http.MethodGet, "/monthlybudget/1", nil)
		ctx = context.WithValue(request.Context(), handlers.UserKey{}, uuid)

		request = request.WithContext(ctx)
		// Manualy inject path
		request.SetPathValue("id", "1")
		fh.GetBudget(response, request)

		var gotBudget data.Budget
		err = json.NewDecoder(response.Body).Decode(&gotBudget)

		if err != nil {
			t.Fatalf("Unable to parse response from server %d into budget, '%v'", response.Body, err)
		}
		want = 1500

		if gotBudget.Budget != want {
			t.Errorf("got %d, want %d", gotBudget.Budget, want)
		}
	})
	t.Run("User cannot update other users budget", func(t *testing.T) {
		b, err := json.Marshal(&data.Budget{
			Budget: 1500,
		})

		if err != nil {
			t.Fatalf("Unable to parse budget from client %d , '%v'", b, err)
		}
		request, _ := http.NewRequest(http.MethodPut, "/monthlybudget/1", bytes.NewReader(b))
		response := httptest.NewRecorder()

		l := slog.Default()
		initDBTestAuth()
		// Manualy inject path
		request.SetPathValue("id", "1")
		fh := handlers.FinanceNewServer(l)

		uuid := TestUserClaims["username"]
		ctx := context.WithValue(request.Context(), handlers.UserKey{}, uuid)

		request = request.WithContext(ctx)
		fh.UpdateBudget(response, request)

		want := 500
		got := response.Code
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

}
