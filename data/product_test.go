package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "nics",
		Price: 2.99,
		SKU:   "asd-asd-asd",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
