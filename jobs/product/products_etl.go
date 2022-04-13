package product

import (
	"fmt"
	"log"
	"time"

	"github.com/juanjoss/off_etl/db"
	"github.com/juanjoss/off_etl/model"
)

type Generator func() <-chan model.ProductRes

type Processor func(<-chan model.ProductRes) <-chan model.ProductRes

type Consumer func(<-chan model.ProductRes)

var (
	apiPageNumber int
	numProducts   int
	apiPageSize   int
)

func RunETL() {
	start := time.Now()

	fmt.Println("running ETL...")
	apiPageNumber = 1
	numProducts = 0

	for i := 1; i <= 2; i++ {
		load()(
			transform()(
				extract(uint(apiPageNumber))(),
			),
		)
		apiPageNumber++

		// getting the number of iterations needed to fetch all product pages
		// containing apiPageSize number of products per page
		if apiPageNumber == 2 {
			fmt.Println("count: ", numProducts)
			fmt.Println("page size: ", apiPageSize)

			iterations := numProducts/apiPageSize + 1
			fmt.Println("iterations: ", iterations)
		}
	}

	duration := time.Since(start)
	fmt.Println(duration)
}

// it creates product batches by fetching pages from the API endpoint
func extract(page uint) Generator {
	return func() <-chan model.ProductRes {
		products := make(chan model.ProductRes)

		productsRes, err := Fetch(page)
		if err != nil {
			log.Printf("error fetching: %v", err)
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

		log.Println(len(productsRes.Products))

		return products
	}
}

// it takes product batches and makes transformations over them
func transform() Processor {
	return func(products <-chan model.ProductRes) <-chan model.ProductRes {
		transformedProducts := make(chan model.ProductRes)

		go func() {
			defer close(transformedProducts)

			for product := range products {
				if product.Name == "" {
					product.Name = "Unkown"
				}
				transformedProducts <- product
			}
		}()

		return transformedProducts
	}
}

func load() Consumer {
	return func(products <-chan model.ProductRes) {
		for {
			p, ok := <-products
			if ok {
				db.Get().Create(p)
			} else {
				log.Fatalf("error during product load process: ok = %v", ok)
			}
		}
	}
}
