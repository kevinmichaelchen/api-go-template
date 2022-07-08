package logging

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"log"
)

var Module = fx.Module("logging",
	fx.Provide(NewLogger),
)

func NewLogger() *zap.Logger {
	// TODO configure log options
	logger, err := zap.NewProduction(
		zap.AddCaller(),
	)
	if err != nil {
		log.Fatalf("failed to build zap logger: %v", err)
	}
	return logger
}
