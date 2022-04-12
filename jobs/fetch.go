package jobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/juanjoss/off_etl/product"
)

func Fetch(page uint) (*product.ProductsRes, error) {
	url := fmt.Sprintf("https://world.openfoodfacts.org?json=true&page=%d", page)
	fmt.Println("fetching URL: ", url)

	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("status is not 200")
	}

	var products product.ProductsRes
	err = json.NewDecoder(res.Body).Decode(&products)
	if err != nil {
		return nil, err
	}

	return &products, nil
}
