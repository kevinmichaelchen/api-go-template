package service

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/kevinmichaelchen/api-go-template/internal/idl/coop/drivers/foo/v1beta1"
)

// ConnectWrapper wraps our gRPC service.
type ConnectWrapper struct {
	s *Service
}

func NewConnectWrapper(s *Service) *ConnectWrapper {
	return &ConnectWrapper{s: s}
}

func (c *ConnectWrapper) GetFoo(
	ctx context.Context,
	req *connect.Request[v1beta1.GetFooRequest],
) (*connect.Response[v1beta1.GetFooResponse], error) {
	res, err := c.s.GetFoo(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	out := connect.NewResponse(res)
	out.Header().Set("API-Version", "v1beta1")
	return out, nil
}

func (c *ConnectWrapper) CreateFoo(
	ctx context.Context,
	req *connect.Request[v1beta1.CreateFooRequest],
) (*connect.Response[v1beta1.CreateFooResponse], error) {
	res, err := c.s.CreateFoo(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	out := connect.NewResponse(res)
	out.Header().Set("API-Version", "v1beta1")
	return out, nil
}
