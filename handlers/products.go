// Package classification of Product API
//
// Documentation for Product API
//
//  Schemes: http
//  BasePath: /
//  Version: 1.0.0
//
//  Consumes:
//  - application/json
//
//  Produces:
//  - application/json
//
// swagger:meta

package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"microserver/data"

	"github.com/gorilla/mux"
)

// A list of products in the response
// swagger:response productsResponse
type ProductsResponseWrapper struct {
	// All products in the system
	// in: body
	Body []data.Product
}

// A single products in the response
// swagger:response productResponse
type ProductResponseWrapper struct {
	// All products in the system
	// in: body
	Body data.Product
}

// swagger:parameters updateProduct
type ProductIdParameterWrapper struct {
	// The id of the product to update
	// in: path
	// required: true
	// pattern: [a-zA-Z]{8}
	ID string `json:"id"`
}

// swagger:response noContent
type ProductsNoContent struct{}

// define the struct for the handler
type Products struct {
	l *log.Logger
}

type KeyProduct struct{}

// creates a new handler and returns a pointer to it
func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

// swagger:route GET /products products listProducts
// returns a list of products
//
//  Produces:
//  - application/json
//
//  Responses:
//   200: productsResponse
//
// get products list
func (h *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	w.Header().Set("content-type", "application/json")
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Encoding error", http.StatusInternalServerError)
	}
}

// swagger:route POST /products products createProduct
// returns the created product
//
//  Consumes:
//  - application/json
//
//  Produces:
//  - application/json
//
//  Responses:
//   201: productResponse
//
// add a new product
func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	data.AddProduct(prod)
}

// swagger:route PUT /products/{id} products updateProduct
// updates the product
//
//  Responses:
//   204: noContent
//
func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	prod.ID = id
	err := data.UpdateProduct(id, prod)
	if err == data.ErrorProductNotFound || err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		// you should use a buffered reader in case the payload is too large
		// so that the memory is not filled with the payload itself
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(w, "Could not decode the payload", http.StatusBadRequest)
			return
		}

		err = prod.Validate()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error validating product: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
