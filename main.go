package main

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/juanjoss/off-etl/jobs"
	"github.com/juanjoss/off-etl/repository"
)

func main() {
	repo := repository.NewRepository()

	jobs.RunBrandsETL(repo)

	s := gocron.NewScheduler(time.UTC)
	s.LimitRunsTo(50)

	s.Every(1).Second().Do(jobs.RunProductsETL, repo)

	s.StartBlocking()
}
