package model

type Repository interface {
	CreateSchema()
	DeleteSchema()

	AddProduct(*Product) error

	AddBrand(*Brand) error
	SearchBrand(string) (*Brand, error)

	AddProductBrand(string, []*Brand)

	AddProductNutrientLevels(*NutrientLevels) (uint8, error)
	GetProductNutrientLevelsId(*NutrientLevels) (uint8, error)
}
