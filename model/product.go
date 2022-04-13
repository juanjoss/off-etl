package model

import "gorm.io/gorm"

type ProductRes struct {
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
	Count    int          `json:"count"`
	Products []ProductRes `json:"products"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
}

type Product struct {
	gorm.Model
	Barcode             string
	Name                string
	Quantity            string
	Categories          string
	Brand               Brand
	Packaging           string
	ImageUrl            string
	ImageIngredientsUrl string
	ImageNutritionUrl   string
}

func (p *Product) ToModel() Product {
	return Product{
		Barcode: p.Barcode,
	}
}
