package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	env_config "github.com/jjma22/finance-tracker/internal/config"
	"github.com/jjma22/finance-tracker/internal/data"
	"github.com/jjma22/finance-tracker/internal/database"
	"github.com/jjma22/finance-tracker/internal/handlers"
)

func initExpDBTest() {
	Config := *env_config.LoadConfig("../../.env")

	db_connection := Config.Database

	// Declare logger
	l := slog.Default()

	// Setup database connection
	database.InitDb(l, &db_connection)
}

var TestUuid = "bkjca-12jk0k-azmp"

func TestPostExpenseSuccess(t *testing.T) {
	t.Run("sets expense for user", func(t *testing.T) {
		expense := data.Expense{
			Name:  "Rent",
			Price: 2000,
			SKU:   "abc-bca-abc",
		}

		e, err := json.Marshal(&expense)

		if err != nil {
			t.Fatalf("Unable to parse expense %v, '%v'", expense, err)
		}

		request, _ := http.NewRequest(http.MethodPost, "/expense", bytes.NewReader(e))
		response := httptest.NewRecorder()

		l := slog.Default()
		initExpDBTest()

		fh := handlers.FinanceNewServer(l)

		ctx := context.WithValue(request.Context(), handlers.Keyexpense{}, &expense)
		ctx = context.WithValue(ctx, handlers.UserKey{}, TestUuid)
		request = request.WithContext(ctx)
		fh.AddExpense(response, request)

		want := 201
		got := response.Code

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}

	})
}

func TestGetExpenseReturnsExpense(t *testing.T) {
	t.Run("Get Expense returns expense successfully", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/expense/1", nil)
		response := httptest.NewRecorder()

		l := slog.Default()
		initExpDBTest()

		fh := handlers.FinanceNewServer(l)

		request.SetPathValue("id", "1")
		ctx := context.WithValue(request.Context(), handlers.UserKey{}, TestUuid)
		request = request.WithContext(ctx)

		fh.GetExpense(response, request)

		var expense data.Expense

		err := json.Unmarshal(response.Body.Bytes(), &expense)
		if err != nil {
			t.Fatalf("Could not parse response into expense, %v", err)
		}

		expectedExpense := data.Expense{
			Name:  "Rent",
			Price: 2000,
			SKU:   "abc-bca-abc",
		}

		if expense.Name != expectedExpense.Name {
			t.Errorf("Failed to return correct expense name; got %s, expected %s", expense.Name, expectedExpense.Name)
		}

		if expense.Price != expectedExpense.Price {
			t.Errorf("Failed to return correct expense price; got %f, expected %f", expense.Price, expectedExpense.Price)
		}

		if expense.SKU != expectedExpense.SKU {
			t.Errorf("Failed to return correct expense SKU; got %s, expected %s", expense.SKU, expectedExpense.SKU)
		}
	})
}

func TestGetExpenseCannotGetUndefinedExpense(t *testing.T) {
	t.Run("Testing get expense fails when expense id does not exist", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/expense/2", nil)
		response := httptest.NewRecorder()

		l := slog.Default()
		initExpDBTest()

		fh := handlers.FinanceNewServer(l)

		request.SetPathValue("id", "2")
		ctx := context.WithValue(request.Context(), handlers.UserKey{}, TestUuid)
		request = request.WithContext(ctx)

		fh.GetExpense(response, request)

		returnedError := strings.TrimSpace(response.Body.String())
		expectedError := "Could not retrieve expense"

		if returnedError != expectedError {
			t.Errorf("Expected error %s, got %s", expectedError, returnedError)
		}

		expectedCode := 500
		if response.Code != expectedCode {
			t.Errorf("Expected %d, got %d", expectedCode, response.Code)
		}

	})
}

func TestPutExpense(t *testing.T) {
	t.Run("Test expense price can be updated", func(t *testing.T) {
		newExpense := data.Expense{
			Price: 2200,
		}

		exp, err := json.Marshal(newExpense)
		if err != nil {
			t.Fatalf("Unable to parse updated expense into JSON, %v", err)
		}

		request, _ := http.NewRequest(http.MethodPut, "/expense/1", bytes.NewReader(exp))
		response := httptest.NewRecorder()

		l := slog.Default()
		initExpDBTest()

		fh := handlers.FinanceNewServer(l)

		request.SetPathValue("id", "1")
		ctx := context.WithValue(request.Context(), handlers.UserKey{}, TestUuid)
		request = request.WithContext(ctx)
		fh.UpdateExpense(response, request)

		expectedCode := 200

		if response.Code != expectedCode {
			t.Errorf("Expected %d, got %d", expectedCode, response.Code)
		}

		request, _ = http.NewRequest(http.MethodGet, "/expense/1", nil)
		request.SetPathValue("id", "1")
		ctx = context.WithValue(request.Context(), handlers.UserKey{}, TestUuid)
		request = request.WithContext(ctx)

		fh.GetExpense(response, request)

		expectedCode = 200

		if response.Code != expectedCode {
			t.Errorf("Expected %d, got %d", expectedCode, response.Code)
		}

		var updatedExpense data.Expense

		err = json.Unmarshal(response.Body.Bytes(), &updatedExpense)
		if err != nil {
			t.Fatalf("Unable to parse expense returned from server, err: %v", err)
		}

		expectedExpense := data.Expense{
			Name:  "Rent",
			Price: 2200,
		}

		if expectedExpense.Name != updatedExpense.Name {
			t.Errorf("Expense name unexpectadly updated. Expected %s, got %s", updatedExpense.Name, expectedExpense.Name)
		}

		if expectedExpense.Price != updatedExpense.Price {
			t.Errorf("Expense price not update. Got %f, expected %f", updatedExpense.Price, expectedExpense.Price)
		}

	})

	t.Run("Test expense name can be updated", func(t *testing.T) {
		newExpense := data.Expense{
			Name: "UpdatedRent",
		}

		exp, err := json.Marshal(newExpense)
		if err != nil {
			t.Fatalf("Unable to parse updated expense into JSON, %v", err)
		}

		request, _ := http.NewRequest(http.MethodPut, "/expense/1", bytes.NewReader(exp))
		response := httptest.NewRecorder()

		l := slog.Default()
		initExpDBTest()

		fh := handlers.FinanceNewServer(l)

		request.SetPathValue("id", "1")
		ctx := context.WithValue(request.Context(), handlers.UserKey{}, TestUuid)
		request = request.WithContext(ctx)
		fh.UpdateExpense(response, request)

		expectedCode := 200

		if response.Code != expectedCode {
			t.Errorf("Expected %d, got %d", expectedCode, response.Code)
		}

		request, _ = http.NewRequest(http.MethodGet, "/expense/1", nil)
		request.SetPathValue("id", "1")
		ctx = context.WithValue(request.Context(), handlers.UserKey{}, TestUuid)
		request = request.WithContext(ctx)

		fh.GetExpense(response, request)

		expectedCode = 200

		if response.Code != expectedCode {
			t.Errorf("Expected %d, got %d", expectedCode, response.Code)
		}

		var updatedExpense data.Expense

		err = json.Unmarshal(response.Body.Bytes(), &updatedExpense)
		if err != nil {
			t.Fatalf("Unable to parse expense returned from server, err: %v", err)
		}

		expectedExpense := data.Expense{
			Name:  "UpdatedRent",
			Price: 2200,
		}

		if expectedExpense.Name != updatedExpense.Name {
			t.Errorf("Expense name unexpectadly updated. Expected %s, got %s", updatedExpense.Name, expectedExpense.Name)
		}

		if expectedExpense.Price != updatedExpense.Price {
			t.Errorf("Expense price not update. Got %f, expected %f", updatedExpense.Price, expectedExpense.Price)
		}

	})
}
