package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jjma22/finance-tracker/internal/database"
	"github.com/jjma22/finance-tracker/internal/handlers"
)

func main() {

	// Declare logger
	l := slog.Default()

	// Setup database connection
	database.InitDb(l)

	// Declare handler
	fh := handlers.FinanceNewServer(l)

	// Initialise new ServerMux
	sm := http.NewServeMux()

	sm.Handle("POST /monthlybudget", fh.MiddleWareValidateBudget(http.HandlerFunc(fh.SetBudget)))
	sm.HandleFunc("GET /monthlybudget/{id}", fh.GetBudget)
	sm.HandleFunc("PUT /monthlybudget/{id}", fh.UpdateBudget)

	sm.HandleFunc("GET /expense/total", fh.GetTotalExpense)
	sm.HandleFunc("GET /expense/{id}", fh.GetExpense)
	sm.HandleFunc("GET /expense", fh.GetExpenses)
	sm.Handle("POST /expense", fh.MiddleWareValidateExpense(http.HandlerFunc(fh.AddExpense)))
	sm.HandleFunc("PUT /expense/update/{id}", fh.UpdateExpense)
	sm.HandleFunc("DELETE /expense/delete/{id}", fh.DeleteExpense)

	//Remove 127.0.0.1 when deploying to Docker, causes issues on local firewall without
	serverPort := "127.0.0.1:9090"

	// Initialise server
	s := &http.Server{
		Addr:         serverPort,
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Error(err.Error())
		}
	}()

	l.Info("Listening on port", "port", serverPort)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Info("Recieve termiante, graceful shutdown", "sig", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

}
