package main

import (
	"github.com/juanjoss/off-etl/jobs"
	"github.com/juanjoss/off-etl/repository"
)

func main() {
	repo := repository.NewRepository()

	jobs.RunBrandsETL(repo)
	jobs.RunProductsETL(repo)
}
