package data

import (
	"encoding/json"
	"errors"
	"io"
	"math"
	"math/rand"
	"time"
)

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// implement fromJSON func
func (p *Product) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

// a type for slice of pointer-to-product
type Products []*Product

// implement toJSON func
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes() string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func generateId() string {
	rand.Seed(time.Now().UnixNano())
	id := make([]byte, 8)
	for i := 0; i < 8; i++ {
		id[i] = byte(99 + math.Floor(rand.Float64()*26))
	}
	return string(id)
}

// dummy products list
var productList = Products{
	{
		ID:          "AsitUyxD",
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.49,
		SKU:         "abc123",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          "ErrtausQ",
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "def456",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}

var ErrorProductNotFound = errors.New("product not found")

func findProduct(id string) (int, error) {
	for i, p := range productList {
		if p.ID == id {
			return i, nil
		}
	}
	return -1, ErrorProductNotFound
}

func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) {
	p.ID = generateId()
	productList = append(productList, p)
}

func UpdateProduct(id string, data *Product) error {
	p, err := findProduct(id)
	if err != nil {
		return err
	}
	productList[p] = data
	return nil
}
