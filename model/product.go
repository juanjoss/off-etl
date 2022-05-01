package model

import (
	"fmt"

	"github.com/juanjoss/off_etl/db"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

var MandatoryStateTags = []string{
	"en:product-name-completed",
	"en:brands-completed",
	"en:quantity-completed",
	"en:front-photo-selected",
	"en:nutrition-facts-completed",
	"en:checked",
}

type ProductsRes struct {
	Count    int          `json:"count"`
	Products []ProductRes `json:"products"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
}

type ProductRes struct {
	Barcode         string            `json:"code"`
	Name            string            `json:"product_name"`
	Quantity        string            `json:"quantity"`
	Brands          []string          `json:"brands_tags"`
	ImageUrl        string            `json:"image_url"`
	NutriscoreGrade string            `json:"nutriscore_grade"`
	NutrientLevels  map[string]string `json:"nutrient_levels"`
	StateTags       []string          `json:"states_tags"`
}

type Product struct {
	gorm.Model
	Barcode  string
	Name     string
	Quantity string
	ImageUrl string
	Brands   []Brand `gorm:"many2many:product_brands;"`
}

func (p *ProductRes) ToModel() (*Product, error) {
	product := &Product{
		Barcode:  p.Barcode,
		Name:     p.Name,
		Quantity: p.Quantity,
		ImageUrl: p.ImageUrl,
		Brands:   []Brand{},
	}

	// search product brands
	for _, brandName := range p.Brands {
		b := Brand{}
		if err := db.Get().Where("tag = ?", brandName).First(&b).Error; err != nil {
			product.Brands = append(product.Brands, Brand{Tag: brandName + ":not-found"})
			continue
		}

		product.Brands = append(product.Brands, b)
	}

	return product, nil
}

func (pr *ProductRes) HasMandatoryStateTags() bool {
	for _, st := range MandatoryStateTags {
		if !slices.Contains(pr.StateTags, st) {
			fmt.Println("product ", pr.Name, " doesn't contains state tag ", st)
			return false
		}
	}

	return true
}
