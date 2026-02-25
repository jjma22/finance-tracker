package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jjma22/finance-tracker.git/internal/handlers"
)

func main() {

	l := log.New(os.Stdout, "fin-api,", log.LstdFlags)
	fh := handlers.FinanceNewServer(l)
	sm := http.NewServeMux()
	sm.Handle("/", fh)

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
			l.Fatal(err)
		}
	}()

	l.Printf("Listening on port %v", serverPort)

	sigChan :=make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <- sigChan
	l.Println("Recieve termiante, graceful shutdown", sig)
	
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.	Shutdown(tc)
	
}