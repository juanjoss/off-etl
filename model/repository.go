package model

type Repository interface {
	CreateSchema()
	DeleteSchema()

	AddProduct(*Product) error

	AddBrand(*Brand) error
	SearchBrand(string) (*Brand, error)

	AddProductBrands(string, []*Brand) error

	AddProductNutrientLevels(*NutrientLevels) (uint8, error)
	GetProductNutrientLevelsId(*NutrientLevels) (uint8, error)
}
