package service

import (
	"context"
	"github.com/kevinmichaelchen/api-go-template/internal/idl/coop/drivers/foo/v1beta1"
	"github.com/kevinmichaelchen/api-go-template/internal/service/db"
)

type Service struct {
	dataStore *db.Store
}

func NewService(dataStore *db.Store) *Service {
	return &Service{
		dataStore: dataStore,
	}
}

func (s *Service) CreateFoo(ctx context.Context, r *foov1beta1.CreateFooRequest) (*foov1beta1.CreateFooResponse, error) {
	//err := validate(r, r)
	//if err != nil {
	//	return nil, err
	//}
	return s.dataStore.CreateFoo(ctx, r)
}

func (s *Service) GetFoo(ctx context.Context, r *foov1beta1.GetFooRequest) (*foov1beta1.GetFooResponse, error) {
	//err := validate(r, r)
	//if err != nil {
	//	return nil, err
	//}
	return s.dataStore.GetFoo(ctx, r)
}
