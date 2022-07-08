package service

import (
	"context"
	"github.com/kevinmichaelchen/api-go-template/internal/idl/coop/drivers/foo/v1beta1"
	"github.com/kevinmichaelchen/api-go-template/internal/service/db"
	"github.com/kevinmichaelchen/api-go-template/internal/service/health"
	healthV1 "google.golang.org/grpc/health/grpc_health_v1"
)

type Service struct {
	dataStore *db.Store
}

func NewService(dataStore *db.Store) *Service {
	return &Service{
		dataStore: dataStore,
	}
}

func (s *Service) CreateFoo(ctx context.Context, r *v1beta1.CreateFooRequest) (*v1beta1.CreateFooResponse, error) {
	//err := validate(r, r)
	//if err != nil {
	//	return nil, err
	//}
	return s.dataStore.CreateFoo(ctx, r)
}

func (s *Service) GetFoo(ctx context.Context, r *v1beta1.GetFooRequest) (*v1beta1.GetFooResponse, error) {
	//err := validate(r, r)
	//if err != nil {
	//	return nil, err
	//}
	return s.dataStore.GetFoo(ctx, r)
}

func (s *Service) Check(ctx context.Context, in *healthV1.HealthCheckRequest) (*healthV1.HealthCheckResponse, error) {
	return health.Check(ctx, in)
}

func (s *Service) Watch(in *healthV1.HealthCheckRequest, srv healthV1.Health_WatchServer) error {
	return health.Watch(in, srv)
}
