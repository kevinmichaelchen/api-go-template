package sql

import (
	"context"
	"database/sql"
	"github.com/XSAM/otelsql"
	_ "github.com/lib/pq"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("sql",
	fx.Provide(
		NewDatabase,
	),
)

func NewDatabase(logger *zap.Logger, lc fx.Lifecycle) (*sql.DB, error) {
	// TODO make configurable
	dsn := "postgres://postgres:postgres@localhost:5432/foo?sslmode=disable"
	db, err := otelsql.Open("postgres", dsn,
		otelsql.WithAttributes(
			semconv.DBSystemPostgreSQL,
		),
	)
	lc.Append(fx.Hook{
		// To mitigate the impact of deadlocks in application startup and
		// shutdown, Fx imposes a time limit on OnStart and OnStop hooks. By
		// default, hooks have a total of 15 seconds to complete. Timeouts are
		// passed via Go's usual context.Context.
		OnStart: func(context.Context) error {
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing DB connection.")
			err := db.Close()
			if err != nil {
				logger.Error("Failed to close DB connection", zap.Error(err))
				return err
			}
			logger.Info("Successfully closed DB connection")
			return err
		},
	})
	return db, err
}
