package jobs

import (
	"log"
	"time"

	"github.com/juanjoss/off-etl/model"
	"github.com/juanjoss/off-etl/ports"
)

var (
	apiPageNumber = 1
	numProducts   = 0
	apiPageSize   = 100
	iterations    = 0
)

func RunProductsETL(repo ports.Repository) {
	start := time.Now()

	log.Println("\nrunning products ETL...")

	if iterations <= numProducts/apiPageSize {
		load(repo)(
			transform()(
				extract(apiPageNumber)(),
			),
		)
		apiPageNumber++
		iterations++
	}

	duration := time.Since(start)
	log.Printf("%v\n", duration)
}

// it creates product batches by fetching pages from the API endpoint
func extract(page int) func() <-chan model.ProductRes {
	return func() <-chan model.ProductRes {
		products := make(chan model.ProductRes)

		productsRes, err := FetchProducts(page, apiPageSize)
		if err != nil {
			log.Fatalf("error fetching: %v", err)
		}

		if apiPageNumber == 1 {
			numProducts = productsRes.Count
		}

		go func() {
			defer close(products)
			for _, p := range productsRes.Products {
				if p.HasMandatoryStateTags() {
					products <- p
				}
			}
		}()

		return products
	}
}

// it takes product batches and makes transformations over them
func transform() func(<-chan model.ProductRes) <-chan model.ProductRes {
	return func(products <-chan model.ProductRes) <-chan model.ProductRes {
		transformedProducts := make(chan model.ProductRes)

		go func() {
			defer close(transformedProducts)

			for p := range products {
				transformedProducts <- p
			}
		}()

		return transformedProducts
	}
}

func load(repo ports.Repository) func(<-chan model.ProductRes) {
	return func(products <-chan model.ProductRes) {
		for {
			pr, ok := <-products
			if ok {
				// convert product response to product model
				product, err := pr.ToModel()
				if err != nil {
					log.Printf("error converting to product model: %v", err.Error())
					continue
				}

				// checking for wrong fields in product
				if product.HasWrongFields() {
					log.Printf("skipping product %v with wrong fields", product)
					continue
				}

				// search for nutrient levels
				id, err := repo.GetProductNutrientLevelsId(&pr.NutrientLevels)
				if err != nil {
					id, err = repo.AddProductNutrientLevels(&pr.NutrientLevels)
					if err != nil {
						log.Printf("error inserting nutrient levels for product %v: %v", product.Name, err.Error())
					}
				}

				if id == 0 {
					log.Printf("skipping product with wrong nutrient levels %v", pr.NutrientLevels)
					continue
				}

				product.NutrientLevelsId = id

				// add product
				err = repo.AddProduct(product)
				if err != nil {
					log.Printf("error inserting product %v: %v", product, err.Error())
				}

				// search product brands and add them
				var brands []*model.Brand

				for _, brandName := range pr.Brands {
					if brandName != "" {
						brand, err := repo.SearchBrand(brandName)
						if err != nil {
							log.Printf("unable to find brand with tag = %v", brandName)
							continue
						}

						brands = append(brands, brand)
					}
				}

				err = repo.AddProductBrands(product.Barcode, brands)
				if err != nil {
					log.Printf("error inserting product brands: %v", err.Error())
				}
			} else {
				log.Printf("products load process finished (error = %v)\n", ok)
				return
			}
		}
	}
}
