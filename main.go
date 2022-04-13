package main

import (
	"github.com/juanjoss/off_etl/db"
	"github.com/juanjoss/off_etl/jobs/brand"
	"github.com/juanjoss/off_etl/jobs/product"
)

func main() {
	db.DeleteSchema()
	db.Migrate()
	brand.RunETL()
	product.RunETL()
}
