package product

import (
	"github.com/juanjoss/off_etl/db"
	"github.com/juanjoss/off_etl/jobs/brand"
	"gorm.io/gorm"
)

type ProductRes struct {
	Barcode             string   `json:"code"`
	Name                string   `json:"product_name"`
	Quantity            string   `json:"quantity"`
	Categories          string   `json:"categories"`
	Brands              []string `json:"brands_tags"`
	Packaging           string   `json:"packaging"`
	ImageUrl            string   `json:"image_url"`
	ImageIngredientsUrl string   `json:"image_ingredients_url"`
	ImageNutritionUrl   string   `json:"image_nutrition_url"`
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
	Packaging           string
	ImageUrl            string
	ImageIngredientsUrl string
	ImageNutritionUrl   string
	Brands              []brand.Brand `gorm:"many2many:product_brands;"`
}

func (p *ProductRes) ToModel() (*Product, error) {
	product := &Product{
		Barcode:             p.Barcode,
		Name:                p.Name,
		Quantity:            p.Quantity,
		Categories:          p.Categories,
		Packaging:           p.Packaging,
		ImageUrl:            p.ImageUrl,
		ImageIngredientsUrl: p.ImageIngredientsUrl,
		ImageNutritionUrl:   p.ImageNutritionUrl,
		Brands:              []brand.Brand{},
	}

	for _, brandName := range p.Brands {
		b, err := brand.GetBrand(brandName)
		if err != nil {
			product.Brands = append(product.Brands, brand.Brand{Tag: brandName + ":not-found"})
			continue
		}

		product.Brands = append(product.Brands, b)
	}

	return product, nil
}

func Migrate() {
	db.Get().AutoMigrate(&Product{})
}

func DeleteSchema() {
	db.Get().Exec("DROP TABLE IF EXISTS products CASCADE")
}
