package adapters

import (
	"github.com/juanjoss/off_etl/adapters/postgres"
	"github.com/juanjoss/off_etl/model"
)

type Services struct {
	Repo model.Repository
}

func NewServices() *Services {
	return &Services{
		Repo: postgres.NewRepo(),
	}
}
