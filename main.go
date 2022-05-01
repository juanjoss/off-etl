package main

import (
	"github.com/juanjoss/off_etl/db"
	"github.com/juanjoss/off_etl/jobs"
)

func main() {
	db.DeleteBrandSchema()
	db.DeleteProductSchema()

	jobs.MigrateBrand()
	jobs.MigrateProduct()

	jobs.RunBrandsETL()
	jobs.RunProductsETL()
}
