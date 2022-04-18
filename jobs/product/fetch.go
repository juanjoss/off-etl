package product

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func Fetch(page uint) (*ProductsRes, error) {
	url := fmt.Sprintf("https://world.openfoodfacts.org/api/v2/search?fields=code,product_name,quantity,categories,brands_tags,packaging,image_url,image_ingredients_url,image_nutrition_url&json=true&page=%d", page)
	fmt.Println("\nfetching URL: ", url)

	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("status is not 200")
	}

	var products ProductsRes
	err = json.NewDecoder(res.Body).Decode(&products)
	if err != nil {
		return nil, err
	}

	return &products, nil
}
