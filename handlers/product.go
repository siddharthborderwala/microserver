package handlers

import (
	"log"
	"net/http"

	"microserver/data"
)

// define the struct for the handler
type Product struct {
	l *log.Logger
}

// creates a new handler and returns a pointer to it
func NewProduct(l *log.Logger) *Product {
	return &Product{l}
}

// implement the ServeHTTP func
func (p *Product) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.GetProducts(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// get products list
func (h *Product) GetProducts(w http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Encoding error", http.StatusInternalServerError)
	}
}
