package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"microserver/handlers"
)

func main() {
	// initialize a custom logger
	// normally, we want to output the logs in a file, but for now we are using os.Stdout
	l := log.New(os.Stdout, "api:  ", log.LstdFlags)

	// initialize all the handlers
	helloHandler := handlers.NewHello(l)
	goodbyeHandler := handlers.NewGoodBye(l)

	// create a custom servemux (requests multiplexer)
	sm := http.NewServeMux()
	sm.Handle("/", helloHandler)
	sm.Handle("/goodbye", goodbyeHandler)

	// we create our own server to configure stuff like timeouts
	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	go func() {
		// this is blocking, hence we put it inside a go func
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	l.Println("Server listening on :6969")

	// we create a channel, that pipes the os signals to sigChan
	sigChan := make(chan os.Signal)
	// we add specific listeners
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// the received message will be store in sig
	sig := <-sigChan
	l.Printf("Received %s, gracefully shutting down", sig.String())

	// create a context with timeout, to let the server shutdown gradefully
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
