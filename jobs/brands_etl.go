package jobs

import (
	"fmt"
	"log"
	"time"

	"github.com/juanjoss/off_etl/model"
)

func RunBrandsETL(repo model.Repository) {
	start := time.Now()
	fmt.Println("\nrunning brands ETL...")

	loadBrands(repo)(
		transformBrands()(
			extractBrands()(),
		),
	)

	duration := time.Since(start)
	fmt.Printf("%v\n", duration)
}

func extractBrands() func() <-chan model.BrandRes {
	return func() <-chan model.BrandRes {
		brands := make(chan model.BrandRes)

		brandsRes, err := FetchBrands()
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

func transformBrands() func(brands <-chan model.BrandRes) <-chan model.BrandRes {
	return func(brands <-chan model.BrandRes) <-chan model.BrandRes {
		transformedBrands := make(chan model.BrandRes)

		go func() {
			defer close(transformedBrands)

			for brand := range brands {
				transformedBrands <- brand
			}
		}()

		return transformedBrands
	}
}

func loadBrands(repo model.Repository) func(brands <-chan model.BrandRes) {
	return func(brands <-chan model.BrandRes) {
		for {
			b, ok := <-brands
			if ok {
				model := b.ToModel()
				repo.AddBrand(model)
			} else {
				log.Printf("brands load process finished (error = %v)", ok)
				return
			}
		}
	}
}
