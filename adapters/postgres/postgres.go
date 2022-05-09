package postgres

import (
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
	data, err := os.ReadFile("../../db/migrate_up.sql")
	if err != nil {
		log.Fatalf("error reading migrate_up: %v", err.Error())
	}

	pr.db.MustExec(string(data))
}

func (pr *PostgresRepo) DeleteSchema() {
	data, err := os.ReadFile("../../db/migrate_down.sql")
	if err != nil {
		log.Fatalf("error reading migrate_down: %v", err.Error())
	}

	pr.db.MustExec(string(data))
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
		_, err := pr.db.Exec("INSERT INTO product_brands (barcode, tag) VALUES ($1, $2)", barcode, brand.Tag)
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
