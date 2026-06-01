package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
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
