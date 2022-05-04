package main

import (
	"github.com/juanjoss/off_etl/adapters"
	"github.com/juanjoss/off_etl/jobs"
)

func main() {
	adapterServices := adapters.NewServices()
	adapterServices.Repo.DeleteSchema()
	adapterServices.Repo.CreateSchema()

	jobs.RunBrandsETL(adapterServices.Repo)
	jobs.RunProductsETL(adapterServices.Repo)
}
