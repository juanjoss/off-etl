package product

import (
	"fmt"
	"log"
	"time"

	"github.com/juanjoss/off_etl/db"
)

type Generator func() <-chan ProductRes

type Processor func(<-chan ProductRes) <-chan ProductRes

type Consumer func(<-chan ProductRes)

var (
	apiPageNumber int
	numProducts   int
	apiPageSize   int
	iterations    = 2
)

func RunETL() {
	start := time.Now()

	fmt.Println("running products ETL...")
	apiPageNumber = 1
	numProducts = 0

	for i := 1; i <= iterations; i++ {
		load()(
			transform()(
				extract(uint(apiPageNumber))(),
			),
		)
		apiPageNumber++

		// getting the number of iterations needed to fetch all product pages
		// containing apiPageSize number of products per page
		if apiPageNumber == 2 {
			iterations = numProducts/apiPageSize + 1
			fmt.Println("iterations: ", iterations)
		}

		if i == 10 {
			break
		}
	}

	duration := time.Since(start)
	fmt.Println(duration)
}

// it creates product batches by fetching pages from the API endpoint
func extract(page uint) Generator {
	return func() <-chan ProductRes {
		products := make(chan ProductRes)

		productsRes, err := Fetch(page)
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
			for _, product := range productsRes.Products {
				products <- product
			}
		}()

		return products
	}
}

// it takes product batches and makes transformations over them
func transform() Processor {
	return func(products <-chan ProductRes) <-chan ProductRes {
		transformedProducts := make(chan ProductRes)

		go func() {
			defer close(transformedProducts)

			for product := range products {
				// check for unnamed products
				if product.Name == "" {
					product.Name = "unknown"
				}

				// check for products without brand
				if len(product.Brands) == 0 {
					product.Brands = append(product.Brands, "unbranded")
				}

				transformedProducts <- product
			}
		}()

		return transformedProducts
	}
}

func load() Consumer {
	return func(products <-chan ProductRes) {
		for {
			p, ok := <-products
			if ok {
				model, err := p.ToModel()
				if err != nil {
					log.Printf("error converting to product model: %v", err.Error())
					continue
				}

				db.Get().Create(model)
			} else {
				log.Printf("products load process finished (error = %v)", ok)
				return
			}
		}
	}
}
