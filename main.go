package main

import (
	"github.com/jasonlvhit/gocron"
	"github.com/juanjoss/off-etl/jobs"
	"github.com/juanjoss/off-etl/repository"
)

func main() {
	repo := repository.NewRepository()

	jobs.RunBrandsETL(repo)

	gocron.Every(1).Minute().Do(jobs.RunProductsETL, repo)

	<-gocron.Start()
}
