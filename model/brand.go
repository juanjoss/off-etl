package model

import (
	"gorm.io/gorm"
)

type BrandRes struct {
	Tag string `json:"id"`
}

type BrandsRes struct {
	Brands []BrandRes `json:"tags"`
}

type Brand struct {
	gorm.Model
	Tag string `gorm:"uniqueIndex"`
}

func (br *BrandRes) ToModel() *Brand {
	return &Brand{
		Tag: br.Tag,
	}
}
