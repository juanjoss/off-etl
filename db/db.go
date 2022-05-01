package db

import (
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

func DeleteProductSchema() {
	Get().Exec("DROP TABLE IF EXISTS products CASCADE")
	Get().Exec("DROP TABLE IF EXISTS product_brands CASCADE")
}

func DeleteBrandSchema() {
	Get().Exec("DROP TABLE IF EXISTS brands CASCADE")
}
