package model

type Repository interface {
	CreateSchema()
	DeleteSchema()

	AddProduct(*Product) error

	AddBrand(*Brand) error
	SearchBrand(tag string) (*Brand, error)

	AddProductBrand(ProductBrands)
}
