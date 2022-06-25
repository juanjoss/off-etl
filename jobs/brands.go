package jobs

import (
	"log"
	"time"

	"github.com/juanjoss/off-etl/model"
	"github.com/juanjoss/off-etl/ports"
)

func RunBrandsETL(repo ports.Repository) {
	start := time.Now()
	log.Println("running brands ETL...")

	if !repo.BrandsLoaded() {
		loadBrands(repo)(
			transformBrands()(
				extractBrands()(),
			),
		)
	}

	duration := time.Since(start)
	log.Printf("brands load process finished (duration=%v)", duration)
}

func extractBrands() func() <-chan model.BrandRes {
	return func() <-chan model.BrandRes {
		brands := make(chan model.BrandRes)

		brandsRes, err := FetchBrands()
		if err != nil {
			log.Fatalf("error fetching brands: %v", err)
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

func loadBrands(repo ports.Repository) func(brands <-chan model.BrandRes) {
	return func(brands <-chan model.BrandRes) {
		for {
			b, ok := <-brands
			if ok {
				model := b.ToModel()
				repo.AddBrand(model)
			} else {
				return
			}
		}
	}
}
