package main

import (
	"github.com/kevinmichaelchen/api-go-template/internal/app"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	a := fx.New(
		app.Module,
		fx.WithLogger(
			func(logger *zap.Logger) fxevent.Logger {
				//return fxevent.NopLogger
				return &fxevent.ZapLogger{Logger: logger}
			},
		),
	)
	a.Run()
}
