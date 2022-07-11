//go:build integration

package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/kevinmichaelchen/api-go-template/internal/idl/coop/drivers/foo/v1beta1"
	"github.com/kevinmichaelchen/api-go-template/internal/models"
	_ "github.com/lib/pq"
	"github.com/sethvargo/go-envconfig"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_CreateFoo(t *testing.T) {
	ctx := context.Background()

	// Connect to DB
	db, err := newDB()
	require.NoError(t, err)

	// Run migrations
	//driver, err := postgres.WithInstance(db, &postgres.Config{})
	//m, err := migrate.NewWithDatabaseInstance(
	//	"file:///migrations",
	//	"postgres", driver)
	//m.Up()

	// Clear table
	_, err = models.Foos().DeleteAll(ctx, db)
	require.NoError(t, err)

	// Run test
	s := &Store{db: db}
	req := &v1beta1.CreateFooRequest{
		Name: "Kevin",
	}
	res, err := s.CreateFoo(ctx, req)
	require.NoError(t, err)
	require.NotEmpty(t, res.GetFoo().GetId())
}

type Config struct {
	DBConfig *DBConfig `env:",prefix=DB_"`
}

type DBConfig struct {
	User string `env:"USER,default=postgres"`
	Pass string `env:"PASS,default=postgres"`
	Host string `env:"HOST,default=localhost"`
	Port int    `env:"PORT,default=5432"`
	Name string `env:"NAME,default=footest"`
}

func newDB() (*sql.DB, error) {
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
	return sql.Open("postgres", dsn)
}
