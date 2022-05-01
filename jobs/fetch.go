package jobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/juanjoss/off_etl/model"
)

var (
	productFields = []string{
		"code",
		"product_name",
		"quantity",
		"categories",
		"brands_tags",
		"countries",
		"packaging",
		"image_url",
		"image_ingredients_url",
		"image_nutrition_url",
		"nutriscore_grade",
		"nutrient_levels",
		"states_tags",
	}
)

func joinUrlFields() string {
	return strings.Join(productFields, ",")
}

func FetchProducts(page uint) (*model.ProductsRes, error) {
	url := fmt.Sprintf("https://world.openfoodfacts.org/api/v2/search?%s&json=true&page=%d", joinUrlFields(), page)
	fmt.Println("fetching URL: ", url)

	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("status is not 200")
	}

	var products model.ProductsRes
	err = json.NewDecoder(res.Body).Decode(&products)
	if err != nil {
		return nil, err
	}

	return &products, nil
}

func FetchBrands() (*model.BrandsRes, error) {
	url := "https://world.openfoodfacts.org/brands.json"
	fmt.Println("fetching URL: ", url)

	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("status is not 200")
	}

	var brands model.BrandsRes
	err = json.NewDecoder(res.Body).Decode(&brands)
	if err != nil {
		return nil, err
	}

	return &brands, nil
}
