package handlers

import (
	"context"
	"log"
	"net/http"

	"microserver/data"

	"github.com/gorilla/mux"
)

// define the struct for the handler
type Product struct {
	l *log.Logger
}

type KeyProduct struct{}

// creates a new handler and returns a pointer to it
func NewProduct(l *log.Logger) *Product {
	return &Product{l}
}

// get products list
func (h *Product) GetProducts(w http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Encoding error", http.StatusInternalServerError)
	}
}

// add a new product
func (p *Product) AddProduct(w http.ResponseWriter, r *http.Request) {
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	data.AddProduct(prod)
}

func (p *Product) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	prod.ID = id
	err := data.UpdateProduct(id, prod)
	if err == data.ErrorProductNotFound || err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (p Product) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		// you should use a buffered reader in case the payload is too large
		// so that the memory is not filled with the payload itself
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(w, "Could not decode the payload", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
