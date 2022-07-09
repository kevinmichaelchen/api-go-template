package sql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/XSAM/otelsql"
	_ "github.com/lib/pq"
	"github.com/sethvargo/go-envconfig"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("sql",
	fx.Provide(
		NewDatabase,
	),
)

type Config struct {
	DBConfig *DBConfig `env:",prefix=DB_"`
}

type DBConfig struct {
	User string `env:"USER,default=postgres"`
	Pass string `env:"PASS,default=postgres"`
	Host string `env:"HOST,default=localhost"`
	Port int    `env:"PORT,default=5432"`
	Name string `env:"NAME,default=foo"`
}

func NewDatabase(logger *zap.Logger, lc fx.Lifecycle) (*sql.DB, error) {
	var cfg Config
	err := envconfig.Process(context.Background(), &cfg)
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBConfig.User,
		cfg.DBConfig.Pass,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.Name,
	)
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
