package product

type Product struct {
	Barcode             string `json:"code"`
	Name                string `json:"product_name"`
	Quantity            string `json:"quantity"`
	Categories          string `json:"categories"`
	Brands              string `json:"brands"`
	Packaging           string `json:"packaging"`
	ImageUrl            string `json:"image_url"`
	ImageIngredientsUrl string `json:"image_ingredients_url"`
	ImageNutritionUrl   string `json:"image_nutrition_url"`
}

type ProductsRes struct {
	Count    int       `json:"count"`
	Products []Product `json:"products"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
}
