package jobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/juanjoss/off_etl/model"
)

const (
	countryDomain = "world"
)

var (
	productFields = []string{
		"code",
		"product_name",
		"quantity",
		"brands_tags",
		"image_url",
		"nutrient_levels",
		"nutriments",
		"nutriscore_grade",
		"states_tags",
	}
)

func joinUrlFields() string {
	return strings.Join(productFields, ",")
}

func FetchProducts(page, pageSize int) (*model.ProductsRes, error) {
	url := fmt.Sprintf("https://%s.openfoodfacts.org/api/v2/search?%s&json=true&pageSize=%d&page=%d", countryDomain, joinUrlFields(), pageSize, page)
	log.Printf("\nfetching URL: %s\n", url)

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
	url := fmt.Sprintf("https://%s.openfoodfacts.org/brands.json", countryDomain)
	log.Printf("\nfetching URL: %s\n", url)

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
