package handlers

import (
	"errors"
	"log"
	"net/http"
	"regexp"

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

func getParam(matcher string, r *http.Request) ([]byte, error) {
	re, err := regexp.Compile(matcher)
	if err != nil {
		return nil, err
	}
	match := re.Find([]byte(r.URL.Path))
	if match == nil {
		return nil, errors.New("no match found")
	}
	return match, nil
}

// implement the ServeHTTP func
func (p *Product) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p.getProducts(w, r)
	case http.MethodPost:
		p.addProduct(w, r)
	case http.MethodPut:
		param, err := getParam(`/([a-zA-Z]{8})\b`, r)
		if err != nil {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
		}
		id := string(param[1:])
		p.updateProduct(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// get products list
func (h *Product) getProducts(w http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Encoding error", http.StatusInternalServerError)
	}
}

// add a new product
func (p *Product) addProduct(w http.ResponseWriter, r *http.Request) {
	prod := &data.Product{}
	// you should use a buffered reader in case the payload is too large
	// so that the memory is not filled with the payload itself
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Could not decode the payload", http.StatusBadRequest)
	}
	data.AddProduct(prod)
}

func (p *Product) updateProduct(w http.ResponseWriter, r *http.Request, id string) {
	prod := &data.Product{}
	// you should use a buffered reader in case the payload is too large
	// so that the memory is not filled with the payload itself
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Could not decode the payload", http.StatusBadRequest)
	}
	err = data.UpdateProduct(id, prod)
	if err == data.ErrorProductNotFound || err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}
