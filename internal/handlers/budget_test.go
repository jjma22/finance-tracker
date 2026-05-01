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

	"github.com/jjma22/finance-tracker/internal/data"
	"github.com/jjma22/finance-tracker/internal/database"
	"github.com/jjma22/finance-tracker/internal/handlers"
)

// Using https://quii.gitbook.io/learn-go-with-tests/build-an-application/http-server as source

func TestPOSTBudget(t *testing.T) {
	t.Run("sets budget", func(t *testing.T) {

		b := &data.Budget{
			Budget: 1000,
		}

		en, err := json.Marshal(&data.Budget{
			Budget: 1000,
		})

		if err != nil {
			t.Fatalf("Unable to parse budget from client %d , '%v'", b, err)
		}
		request, _ := http.NewRequest(http.MethodPost, "/monthlybudget", bytes.NewReader(en))
		response := httptest.NewRecorder()
		l := slog.Default()

		database.InitDb(l)
		// Manualy inject path
		request.SetPathValue("id", "1")
		fh := handlers.FinanceNewServer(l)

		// Set context due to middleware
		ctx := context.WithValue(request.Context(), handlers.Budget{}, b)
		request = request.WithContext(ctx)

		fh.SetBudget(response, request)

		want := 200
		got := response.Code

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}

	})
}

func TestGETBudget(t *testing.T) {
	t.Run("returns budget", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/monthlybudget/1", nil)
		response := httptest.NewRecorder()
		l := slog.Default()
		// may be worth looking https://golang.testcontainers.org/
		database.InitDb(l)
		// Manualy inject path
		request.SetPathValue("id", "1")
		fh := handlers.FinanceNewServer(l)
		fh.GetBudget(response, request)

		var got data.Budget
		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse response from server %d into budgret, '%v'", response.Body, err)
		}
		want := 1000

		if got.Budget != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}

func TestPUTBudget(t *testing.T) {
	t.Run("Updates budget", func(t *testing.T) {

		b, err := json.Marshal(&data.Budget{
			Budget: 1500,
		})

		if err != nil {
			t.Fatalf("Unable to parse budget from client %d , '%v'", b, err)
		}
		request, _ := http.NewRequest(http.MethodPut, "/monthlybudget/1", bytes.NewReader(b))
		response := httptest.NewRecorder()
		l := slog.Default()
		// may be worth looking https://golang.testcontainers.org/
		database.InitDb(l)
		// Manualy inject path
		request.SetPathValue("id", "1")
		fh := handlers.FinanceNewServer(l)
		fh.UpdateBudget(response, request)

		want := 200
		got := response.Code
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}

		// Check id 1 budget has updated
		request, _ = http.NewRequest(http.MethodGet, "/monthlybudget/1", nil)
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

}

// Copilot example
func TestBudgetFromJSON(t *testing.T) {
	body := strings.NewReader(`{"budget":1000}`)
	req := httptest.NewRequest(http.MethodPost, "/monthlybudget", body)

	l := slog.Default()
	fh := handlers.FinanceNewServer(l)

	v, err := fh.BudgetFromJSON(req)
	if err != nil {
		t.Fatal(err)
	}
	if v.Budget != 1000 {
		t.Errorf("got %d, want %d", v.Budget, 1000)
	}
}

func TestBadBudgetFromJSON(t *testing.T) {
	body := strings.NewReader(`{budget:1000}`)
	req := httptest.NewRequest(http.MethodPost, "/monthlybudget", body)

	l := slog.Default()
	fh := handlers.FinanceNewServer(l)

	_, err := fh.BudgetFromJSON(req)
	if err == nil {
		t.Fatal("No error when parsing bad JSON")
	}
}

// Validate middleware stops bad budgets
func TestMiddleWareValidateBudget(t *testing.T) {
	body := strings.NewReader(`{"budget":-1000}`)
	req := httptest.NewRequest(http.MethodPost, "/monthlybudget", body)

	l := slog.Default()
	fh := handlers.FinanceNewServer(l)

	response := httptest.NewRecorder()

	rsp := fh.MiddleWareValidateBudget(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	}))

	rsp.ServeHTTP(response, req)

	badResponse := 400

	if response.Code != badResponse {
		t.Errorf("got %d, wanted %d", response.Code, badResponse)
	}

}

// Test middleware correctly passes on request body
func TestMiddleWarePassesOnBudget(t *testing.T) {
	b, err := json.Marshal(&data.Budget{
		Budget: 1000,
	})

	if err != nil {
		t.Fatalf("Unable to parse budget from client %d , '%v'", b, err)
	}
	req, _ := http.NewRequest(http.MethodPost, "/monthlybudget", bytes.NewReader(b))
	//req := httptest.NewRequest(http.MethodPost, "/monthlybudget/1", body)

	l := slog.Default()
	fh := handlers.FinanceNewServer(l)

	response := httptest.NewRecorder()

	f1 := func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("About to parse")

		b := r.Context().Value(handlers.Budget{}).(*data.Budget)

		if b.Budget != 1000 {
			t.Errorf("Budget was not successfully passed to next handler")
		}
	}
	rsp := fh.MiddleWareValidateBudget(http.HandlerFunc(f1))

	rsp.ServeHTTP(response, req)
}
