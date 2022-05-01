package jobs

import (
	"github.com/juanjoss/off_etl/db"
	"github.com/juanjoss/off_etl/model"
)

func MigrateProduct() {
	db.Get().AutoMigrate(&model.Product{})
}

func MigrateBrand() {
	db.Get().AutoMigrate(&model.Brand{})
}
