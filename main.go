package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"microserver/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// seed random
	rand.Seed(time.Now().UnixNano())
	// initialize a custom logger
	// normally, we want to output the logs in a file, but for now we are using os.Stdout
	logger := log.New(os.Stdout, "API:  ", log.LstdFlags)

	// initialize all the handlers
	productHandler := handlers.NewProduct(logger)

	// create a custom servemux (requests multiplexer)
	sm := mux.NewRouter()

	// create a get router
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", productHandler.GetProducts)

	// create a put router
	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.Use(productHandler.MiddlewareProductValidation)
	putRouter.HandleFunc("/products/{id:[a-zA-Z]{8}}", productHandler.UpdateProduct)

	// create a post router
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.Use(productHandler.MiddlewareProductValidation)
	postRouter.HandleFunc("/products", productHandler.AddProduct)

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
			logger.Fatal(err)
		}
	}()

	logger.Println("Server listening on :9090")

	// we create a channel, that pipes the os signals to sigChan
	sigChan := make(chan os.Signal)
	// we add specific listeners
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// the received message will be store in sig
	sig := <-sigChan
	logger.Printf("Received %s, gracefully shutting down", sig.String())

	// create a context with timeout, to let the server shutdown gradefully
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
