package db

import (
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/juanjoss/off_etl/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Get() *gorm.DB {
	if DB == nil {
		dsn := "host=localhost user=root password=root dbname=off_etl port=5432 sslmode=disable"
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("unable to connect to DB: %v", err.Error())
		}

		DB = db
	}

	return DB
}

func Migrate() {
	Get().AutoMigrate(&model.Brand{}, &model.Product{})
}

func DeleteSchema() {
	Get().Exec("DELETE FROM brands")
	Get().Exec("DELETE FROM products")
}
