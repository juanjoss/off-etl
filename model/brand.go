package model

import "gorm.io/gorm"

type BrandRes struct {
	Name string `json:"name"`
}

type BrandsRes struct {
	Brands []BrandRes `json:"tags"`
}

type Brand struct {
	gorm.Model
	Name string
}

func (br *BrandRes) ToModel() *Brand {
	return &Brand{
		Name: br.Name,
	}
}
