package postgres

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/juanjoss/off_etl/model"
	_ "github.com/lib/pq"
)

type PostgresRepo struct {
	db *sqlx.DB
}

func NewRepo() *PostgresRepo {
	db, err := sqlx.Connect("postgres", "user=root password=root dbname=off_etl sslmode=disable")
	if err != nil {
		log.Fatalf("unable to connect to DB: %v", err.Error())
	}

	repo := &PostgresRepo{
		db: db,
	}

	return repo
}

func (pr *PostgresRepo) CreateSchema() {
	data, err := os.ReadFile("./db/init.sql")
	if err != nil {
		log.Fatalf("error reading init.sql: %v", err.Error())
	}

	pr.db.MustExec(string(data))
}

func (pr *PostgresRepo) DeleteSchema() {
	pr.db.MustExec(`
		DROP TABLE IF EXISTS 
			products, 
			product_brands, 
			brands, 
			nutritional_information 
		CASCADE
	`)
}

/*
	Product
*/
func (pr *PostgresRepo) AddProduct(p *model.Product) error {
	_, err := pr.db.NamedExec(`
		INSERT INTO products (barcode, name, quantity, image_url)
		VALUES (:barcode, :name, :quantity, :image_url)`, p)
	if err != nil {
		return err
	}

	return nil
}

/*
	Brand
*/
func (pr *PostgresRepo) AddBrand(b *model.Brand) error {
	_, err := pr.db.NamedExec(`
		INSERT INTO brands (tag)
		VALUES (:tag)`, b)
	if err != nil {
		return err
	}

	return nil
}

func (pr *PostgresRepo) SearchBrand(tag string) (*model.Brand, error) {
	b := &model.Brand{}
	err := pr.db.Get(b, "SELECT tag FROM brands WHERE tag=$1", tag)
	if err != nil {
		return nil, nil
	}

	return b, nil
}

/*
	Product Brands
*/
func (pr *PostgresRepo) AddProductBrand(productBrands model.ProductBrands) {
	pr.db.NamedExec("INSERT INTO products (barcode, name, quantity, image_url) VALUES (:barcode, :name, :quantity, :image_url)", productBrands.Product)

	for _, brand := range productBrands.Brands {
		fmt.Println("inserting barcode: ", productBrands.Product.Barcode, ", tag: ", brand.Tag)

		pr.db.MustExec("INSERT INTO product_brands (barcode, tag) VALUES ($1, $2)", productBrands.Product.Barcode, brand.Tag)
	}
}
