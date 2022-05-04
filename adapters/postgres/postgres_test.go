package postgres

import (
	"testing"

	"github.com/juanjoss/off_etl/model"
)

func TestPostgresRepo(t *testing.T) {
	pr := NewRepo()
	pr.CreateSchema()

	p := &model.Product{
		Barcode:  "5449000131805",
		Name:     "Coca-Cola Zero",
		Quantity: "0.33 l",
		ImageUrl: "https://images.openfoodfacts.org/images/products/544/900/013/1805/front_en.490.400.jpg",
	}

	err := pr.AddProduct(p)
	if err != nil {
		t.Fail()
	}

	pr.DeleteSchema()
}
