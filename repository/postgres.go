package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/juanjoss/off-etl/model"
	_ "github.com/lib/pq"
)

type PostgresRepo struct {
	db *sqlx.DB
}

func NewRepository() *PostgresRepo {
	host := os.Getenv("DB_HOST")
	driver := os.Getenv("DB_DRIVER")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("SSL_MODE")

	source := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, dbPort, dbUser, dbPassword, dbName, sslMode)

	db, err := sqlx.Connect(driver, source)
	if err != nil {
		log.Fatalf("unable to connect to DB: %v", err.Error())
	}

	repo := &PostgresRepo{
		db: db,
	}

	return repo
}

/*
	Product
*/
func (pr *PostgresRepo) AddProduct(product *model.Product) error {
	_, err := pr.db.NamedExec(`
		INSERT INTO products (
			barcode,
			name,
			quantity,
			image_url,
			energy_100g,
			energy_serving,
			nutrient_levels_id,
			nova_group,
			nutriscore_score,
			nutriscore_grade
		)
		VALUES (:barcode, 
			:name, 
			:quantity, 
			:image_url,
			:energy_100g,
			:energy_serving,
			:nutrient_levels_id,
			:nova_group,
			:nutriscore_score,
			:nutriscore_grade
		)`,
		product,
	)
	if err != nil {
		return err
	}

	return nil
}

/*
	Brand
*/
func (pr *PostgresRepo) AddBrand(brand *model.Brand) error {
	_, err := pr.db.NamedExec("INSERT INTO brands (tag) VALUES (:tag)", brand)
	if err != nil {
		return err
	}

	return nil
}

func (pr *PostgresRepo) SearchBrand(tag string) (*model.Brand, error) {
	b := &model.Brand{}
	err := pr.db.Get(b, "SELECT tag FROM brands WHERE tag=$1", tag)
	if err != nil {
		return b, nil
	}

	return b, nil
}

/*
	Product Brands
*/
func (pr *PostgresRepo) AddProductBrands(barcode string, brands []*model.Brand) error {
	for _, brand := range brands {
		_, err := pr.db.Exec(`
			INSERT INTO product_brands (barcode, tag)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING`,
			barcode, brand.Tag,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
	Product Nutrient Levels
*/
func (pr *PostgresRepo) AddProductNutrientLevels(nl *model.NutrientLevels) (uint8, error) {
	var id uint8
	row := pr.db.QueryRow(`
		INSERT INTO nutrient_levels (
			fat, 
			saturated_fat, 
			sugar, 
			salt
		) 
		VALUES (
			$1, 
			$2, 
			$3, 
			$4
		)
		RETURNING id`,
		nl.Fat,
		nl.SaturatedFat,
		nl.Sugar,
		nl.Salt,
	)
	if err := row.Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

func (pr *PostgresRepo) GetProductNutrientLevelsId(nl *model.NutrientLevels) (uint8, error) {
	var id uint8
	err := pr.db.Get(&id, `
		SELECT id 
		FROM nutrient_levels 
		WHERE 
			fat = $1 AND 
			saturated_fat = $2 AND 
			sugar = $3 AND 
			salt = $4`,
		nl.Fat,
		nl.SaturatedFat,
		nl.Sugar,
		nl.Salt,
	)
	if err != nil {
		return id, err
	}

	return id, nil
}

/*
	ProductOrder event generator function
*/
func (pr *PostgresRepo) GetRandomProductFromUserSsd() (int, int, string, error) {
	var userId, ssdId int
	var barcode string

	rows, err := pr.db.Query(`
		SELECT users.id AS userId, ssds.id AS ssdId, product_ssds.barcode AS barcode
		FROM (users JOIN ssds ON users.id = ssds.id) JOIN product_ssds ON ssds.id = product_ssds.ssd_id
		ORDER BY RANDOM() 
		LIMIT 1
	`)
	if err != nil {
		log.Printf("error querying: %v", err.Error())
	}

	for rows.Next() {
		err = rows.Scan(&userId, &ssdId, &barcode)
		if err != nil {
			log.Printf("error scanning rows: %v", err.Error())
		}
	}
	if err = rows.Err(); err != nil {
		return userId, ssdId, barcode, err
	}

	return userId, ssdId, barcode, nil
}
