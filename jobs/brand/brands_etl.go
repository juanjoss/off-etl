package brand

import (
	"fmt"
	"log"
	"time"

	"github.com/juanjoss/off_etl/db"
)

type Generator func() <-chan BrandRes

type Processor func(<-chan BrandRes) <-chan BrandRes

type Consumer func(<-chan BrandRes)

func RunETL() {
	start := time.Now()

	fmt.Println("running brands ETL...")

	load()(
		transform()(
			extract()(),
		),
	)

	duration := time.Since(start)
	fmt.Println(duration)
}

// it creates product batches by fetching pages from the API endpoint
func extract() Generator {
	return func() <-chan BrandRes {
		brands := make(chan BrandRes)

		brandsRes, err := Fetch()
		if err != nil {
			log.Fatalf("error fetching: %v", err)
		}

		go func() {
			defer close(brands)
			for _, brand := range brandsRes.Brands {
				brands <- brand
			}
		}()

		return brands
	}
}

// it takes product batches and makes transformations over them
func transform() Processor {
	return func(brands <-chan BrandRes) <-chan BrandRes {
		transformedBrands := make(chan BrandRes)

		go func() {
			defer close(transformedBrands)

			for brand := range brands {
				transformedBrands <- brand
			}
		}()

		return transformedBrands
	}
}

func load() Consumer {
	return func(brands <-chan BrandRes) {
		for {
			b, ok := <-brands
			if ok {
				db.Get().Create(b.ToModel())
			} else {
				log.Printf("brands load process finished (error = %v)", ok)
				return
			}
		}
	}
}
