package main

import (
	"github.com/juanjoss/off_etl/jobs/brand"
	"github.com/juanjoss/off_etl/jobs/product"
)

func main() {
	brand.DeleteSchema()
	product.DeleteSchema()

	brand.Migrate()
	product.Migrate()

	brand.RunETL()
	product.RunETL()
}
