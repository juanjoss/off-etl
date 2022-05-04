package model

type BrandRes struct {
	Tag string `json:"id"`
}

type BrandsRes struct {
	Brands []BrandRes `json:"tags"`
}

type Brand struct {
	Tag string `db:"tag"`
}

func (br *BrandRes) ToModel() *Brand {
	return &Brand{
		Tag: br.Tag,
	}
}
