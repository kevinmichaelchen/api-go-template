package service

import (
	"database/sql"
	"github.com/kevinmichaelchen/api-go-template/internal/service"
	"github.com/kevinmichaelchen/api-go-template/internal/service/db"
	"go.uber.org/fx"
)

var Module = fx.Module("service",
	fx.Provide(
		NewService,
		NewDataStore,
	),
)

type ServiceParams struct {
	fx.In
	DataStore *db.Store
}

func NewService(p ServiceParams) *service.Service {
	return service.NewService(p.DataStore)
}

func NewDataStore(sqlDB *sql.DB) *db.Store {
	return db.NewStore(sqlDB)
}
