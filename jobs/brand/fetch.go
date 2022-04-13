package brand

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/juanjoss/off_etl/model"
)

func Fetch() (*model.BrandsRes, error) {
	url := fmt.Sprintf("https://world.openfoodfacts.org/brands.json")
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
