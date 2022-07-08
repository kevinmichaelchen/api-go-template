package db

import (
	"context"
	"github.com/kevinmichaelchen/api-go-template/internal/idl/coop/drivers/foo/v1beta1"
	"github.com/kevinmichaelchen/api-go-template/internal/models"
	"github.com/rs/xid"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (s *Store) CreateFoo(ctx context.Context, r *v1beta1.CreateFooRequest) (*v1beta1.CreateFooResponse, error) {
	foo := &models.Foo{
		ID:   xid.New().String(),
		Name: r.GetName(),
	}
	err := foo.Insert(ctx, s.db, boil.Infer())
	if err != nil {
		return nil, err
	}
	return &v1beta1.CreateFooResponse{
		Foo: fooModelToProto(foo),
	}, nil
}

func (s *Store) GetFoo(ctx context.Context, r *v1beta1.GetFooRequest) (*v1beta1.GetFooResponse, error) {
	foo, err := models.FindFoo(ctx, s.db, r.GetId())
	if err != nil {
		return nil, err
	}
	return &v1beta1.GetFooResponse{
		Foo: fooModelToProto(foo),
	}, nil
}

func fooModelToProto(foo *models.Foo) *v1beta1.Foo {
	return &v1beta1.Foo{
		Id:   foo.ID,
		Name: foo.Name,
	}
}
