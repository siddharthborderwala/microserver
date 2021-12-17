package data

import (
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"required,gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// validate method on *Product
func (p *Product) Validate() error {
	validate := validator.New()
	// custom validator
	validate.RegisterValidation("sku", ValidateSKU)

	return validate.Struct(p)
}

func ValidateSKU(fl validator.FieldLevel) bool {
	re := regexp.MustCompile("^(([a-zA-Z]{3}-){2}[a-zA-Z]{3})$")
	return re.FindString(fl.Field().String()) != ""
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

const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateId() string {
	b := make([]byte, 8)
	size := len(characters)
	for i := range b {
		b[i] = characters[rand.Intn(size)]
	}
	return string(b)
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
