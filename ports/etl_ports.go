package ports

import (
	"github.com/juanjoss/off_etl/model"
)

type Repository interface {
	AddProduct(*model.Product) error

	AddBrand(*model.Brand) error
	SearchBrand(string) (*model.Brand, error)

	AddProductBrands(string, []*model.Brand) error

	AddProductNutrientLevels(*model.NutrientLevels) (uint8, error)
	GetProductNutrientLevelsId(*model.NutrientLevels) (uint8, error)

	GetRandomProductFromUserSsd() (int, int, string, error)
}
