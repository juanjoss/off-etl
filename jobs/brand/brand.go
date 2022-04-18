package brand

import (
	"github.com/juanjoss/off_etl/db"
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
	Tag string
}

func (br *BrandRes) ToModel() *Brand {
	return &Brand{
		Tag: br.Tag,
	}
}

func Migrate() {
	db.Get().AutoMigrate(&Brand{})
	db.Get().Create(&Brand{Tag: "unbranded"})
}

func DeleteSchema() {
	db.Get().Exec("DROP TABLE IF EXISTS brands CASCADE")
}

func GetBrand(name string) (Brand, error) {
	b := Brand{}
	if err := db.Get().Where("tag = ?", name).First(&b).Error; err != nil {
		return Brand{}, err
	}

	return b, nil
}
