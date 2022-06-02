package main

import (
	"github.com/jasonlvhit/gocron"
	"github.com/juanjoss/off_etl/jobs"
	"github.com/juanjoss/off_etl/repository"
)

func main() {
	repo := repository.NewRepository()

	jobs.RunBrandsETL(repo)

	gocron.Every(1).Minute().Do(jobs.RunProductsETL, repo)

	<-gocron.Start()
}
