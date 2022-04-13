package brand

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/juanjoss/off_etl/db"
	"github.com/juanjoss/off_etl/model"
)

type Generator func() <-chan model.BrandRes

type Processor func(<-chan model.BrandRes) <-chan model.BrandRes

type Consumer func(<-chan model.BrandRes)

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
	return func() <-chan model.BrandRes {
		brands := make(chan model.BrandRes)

		brandsRes, err := Fetch()
		if err != nil {
			log.Printf("error fetching: %v", err)
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
	return func(brands <-chan model.BrandRes) <-chan model.BrandRes {
		transformedBrands := make(chan model.BrandRes)

		go func() {
			defer close(transformedBrands)

			for brand := range brands {
				// replacing "-" by spaces in brand name
				if strings.Contains(brand.Name, "-") {
					brand.Name = strings.ReplaceAll(brand.Name, "-", " ")
				}

				transformedBrands <- brand
			}
		}()

		return transformedBrands
	}
}

func load() Consumer {
	return func(brands <-chan model.BrandRes) {
		for {
			b, ok := <-brands
			if ok {
				db.Get().Create(b)
			} else {
				return
			}
		}
	}
}
