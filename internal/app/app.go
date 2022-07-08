package app

import (
	"github.com/kevinmichaelchen/api-go-template/internal/app/config"
	"github.com/kevinmichaelchen/api-go-template/internal/app/grpc"
	"github.com/kevinmichaelchen/api-go-template/internal/app/logging"
	"github.com/kevinmichaelchen/api-go-template/internal/app/metrics"
	"github.com/kevinmichaelchen/api-go-template/internal/app/service"
	"github.com/kevinmichaelchen/api-go-template/internal/app/sql"
	"github.com/kevinmichaelchen/api-go-template/internal/app/tracing"
	"go.uber.org/fx"
)

var Module = fx.Options(
	config.Module,
	grpc.Module,
	logging.Module,
	metrics.Module,
	service.Module,
	sql.Module,
	tracing.Module,
)
