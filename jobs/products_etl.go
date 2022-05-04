package jobs

import (
	"fmt"
	"log"
	"time"

	"github.com/juanjoss/off_etl/model"
)

var (
	apiPageNumber int
	numProducts   int
	apiPageSize   int
	iterations    = 2
)

func RunProductsETL(repo model.Repository) {
	start := time.Now()

	fmt.Println("\nrunning products ETL...")
	apiPageNumber = 1
	numProducts = 0

	for i := 1; i <= iterations; i++ {
		load(repo)(
			transform()(
				extract(uint(apiPageNumber))(),
			),
		)
		apiPageNumber++

		// getting the number of iterations needed to fetch all product pages
		// containing apiPageSize number of products per page
		if apiPageNumber == 2 {
			iterations = numProducts/apiPageSize + 1
		}

		if i == 1 {
			break
		}
	}

	duration := time.Since(start)
	fmt.Printf("%v\n", duration)
}

// it creates product batches by fetching pages from the API endpoint
func extract(page uint) func() <-chan model.ProductRes {
	return func() <-chan model.ProductRes {
		products := make(chan model.ProductRes)

		productsRes, err := FetchProducts(page)
		if err != nil {
			log.Fatalf("error fetching: %v", err)
		}

		if page == 1 {
			apiPageNumber = productsRes.Page
			numProducts = productsRes.Count
			apiPageSize = productsRes.PageSize
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

func load(repo model.Repository) func(<-chan model.ProductRes) {
	return func(products <-chan model.ProductRes) {
		for {
			pr, ok := <-products
			if ok {
				var productBrands model.ProductBrands

				product, err := pr.ToModel()
				if err != nil {
					log.Printf("error converting to product model: %v", err.Error())
					continue
				}

				productBrands.Product = product

				// search product brands
				for _, brandName := range pr.Brands {
					brand, err := repo.SearchBrand(brandName)
					if err != nil {
						log.Printf("unable to find brand with tag = %v", brandName)
						continue
					}

					productBrands.Brands = append(productBrands.Brands, brand)
				}

				repo.AddProductBrand(productBrands)
			} else {
				log.Printf("products load process finished (error = %v)", ok)
				return
			}
		}
	}
}
