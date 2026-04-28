package handlers_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jjma22/finance-tracker/internal/database"
	"github.com/jjma22/finance-tracker/internal/handlers"
)

// Using https://quii.gitbook.io/learn-go-with-tests/build-an-application/http-server as source

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

		got := response.Body.String()
		want := "{\"budget\":2000}"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
