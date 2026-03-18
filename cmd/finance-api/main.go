package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jjma22/finance-tracker.git/internal/database"
	"github.com/jjma22/finance-tracker.git/internal/handlers"
)


func main() {

	//l := log.New(os.Stdout, "fin-api,", log.LstdFlags)
	l := slog.Default()

	// Setup database connection
	database.InitDb(l)

	fh := handlers.FinanceNewServer(l)
	sm := http.NewServeMux()
	//sm.Handle("/", fh)
	//HandleFunc(pattern string, handler func(ResponseWriter, *Request))
	sm.HandleFunc("GET /monthlybudget", fh.GetBudget)
	sm.HandleFunc("PUT /monthlybudget", fh.UpdateBudget)

	sm.HandleFunc("GET /expense", fh.GetExpenses)
	sm.HandleFunc("GET /expense/total", fh.GetTotalExpense)
	sm.Handle("POST /expense", fh.MiddleWareValidateExpense(http.HandlerFunc(fh.AddExpense)))
	sm.HandleFunc("PUT /expense/update/{id}", fh.UpdateExpense)
	sm.HandleFunc("DELETE /expense/delete/{id}", fh.DeleteExpense)


	//Remove 127.0.0.1 when deploying to Docker, causes issues on local firewall without
	serverPort := "127.0.0.1:9090"

	s := &http.Server{
		Addr: serverPort,
		Handler: sm,
		IdleTimeout: 120*time.Second,
		ReadTimeout: 1*time.Second,
		WriteTimeout: 1*time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Error(err.Error())
		}
	}()

	l.Info("Listening on port", "port", serverPort)

	sigChan :=make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <- sigChan
	l.Info("Recieve termiante, graceful shutdown", "sig", sig)
	
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.	Shutdown(tc)
	
}