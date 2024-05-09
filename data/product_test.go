package data

import "testing"

func TestCheckValidation(t *testing.T) {
	p := &Product{
		Name: "ChaiTea",
		Price: 2.50,
		SKU: "qwe-asd-zxc",
	}
	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}