package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jjma22/finance-tracker/internal/auth"
	env_config "github.com/jjma22/finance-tracker/internal/config"
	"github.com/jjma22/finance-tracker/internal/database"
	"github.com/jjma22/finance-tracker/internal/handlers"
)

func main() {

	Config := *env_config.LoadConfig()

	db_connection := Config.Database

	// Declare logger
	l := slog.Default()

	// Setup database connection
	database.InitDb(l, &db_connection)

	// Set up JWT
	auth.InitjwtKey(&Config.Auth)

	// Declare handler
	fh := handlers.FinanceNewServer(l)

	// Initialise new ServerMux
	sm := http.NewServeMux()

	sm.HandleFunc("POST /login", fh.VerifyUser)
	sm.HandleFunc("POST /create/user", fh.CreateUser)

	sm.Handle("POST /monthlybudget", fh.MiddleWareValidateAuthentication(fh.MiddleWareValidateBudget(http.HandlerFunc(fh.SetBudget))))
	sm.Handle("GET /monthlybudget/{id}", fh.MiddleWareValidateAuthentication(http.HandlerFunc(fh.GetBudget)))
	sm.Handle("PUT /monthlybudget/{id}", fh.MiddleWareValidateAuthentication(http.HandlerFunc(fh.UpdateBudget)))

	sm.Handle("GET /expense/total", fh.MiddleWareValidateAuthentication(http.HandlerFunc(fh.GetTotalExpense)))
	sm.Handle("GET /expense/{id}", fh.MiddleWareValidateAuthentication(http.HandlerFunc(fh.GetExpense)))
	sm.Handle("GET /expense", fh.MiddleWareValidateAuthentication(http.HandlerFunc((fh.GetExpenses))))
	sm.Handle("POST /expense", fh.MiddleWareValidateAuthentication(fh.MiddleWareValidateExpense(http.HandlerFunc(fh.AddExpense))))
	sm.Handle("PUT /expense/update/{id}", fh.MiddleWareValidateAuthentication(http.HandlerFunc(fh.UpdateExpense)))
	sm.Handle("DELETE /expense/delete/{id}", fh.MiddleWareValidateAuthentication(http.HandlerFunc(fh.DeleteExpense)))

	//Remove 127.0.0.1 when deploying to Docker, causes issues on local firewall without
	serverPort := ":" + Config.Main.Port

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
