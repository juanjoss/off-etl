package jobs

import (
	"fmt"
	"log"
	"time"

	"github.com/juanjoss/off_etl/product"
)

var (
	apiPageNumber int
	numProducts   int
	apiPageSize   int
)

type Generator func() <-chan product.Product

type Processor func(<-chan product.Product) <-chan product.Product

type Consumer func(<-chan product.Product)

func ETL() {
	start := time.Now()

	fmt.Println("running ETL...")
	apiPageNumber = 1
	numProducts = 0

	for i := 1; i <= 2; i++ {
		printConsumer()(
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
	return func() <-chan product.Product {
		products := make(chan product.Product)

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

		return products
	}
}

// streaming generator
// func randNumGenerator(max int) Generator {
// 	return func() <-chan int {
// 		out := make(chan int, 10)
// 		rand.Seed(time.Now().UnixNano())
// 		go func() {
// 			for {
// 				out <- rand.Intn(max)
// 				time.Sleep(10 * time.Millisecond)
// 			}
// 		}()
// 		return out
// 	}
// }

// it takes product batches and makes transformations over them
func transform() Processor {
	return func(products <-chan product.Product) <-chan product.Product {
		transformedProducts := make(chan product.Product)

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

// it makes use of the transformed products
func printConsumer() Consumer {
	return func(products <-chan product.Product) {
		for {
			p, ok := <-products
			if ok {
				log.Println(p.Barcode, p.Name)
			} else {
				return
			}
		}
	}
}
