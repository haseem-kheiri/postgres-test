package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresConfig struct {
	Image    string
	DB       string
	User     string
	Password string
}

func SetupPostgres(
	t testing.TB,
	cfg PostgresConfig) string {
	t.Helper()

	image := cfg.Image
	if image == "" {
		image = "postgres:16-alpine"
	}

	ctx := context.Background()

	pgContainer, err := postgres.Run(
		ctx,
		image,
		postgres.WithDatabase(cfg.DB),
		postgres.WithUsername(cfg.User),
		postgres.WithPassword(cfg.Password),
		testcontainers.WithWaitStrategy(
			// wait.ForLog returns a LogStrategy pointer which
			// HAS the methods WithOccurrence and WithStartupTimeout
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(1).
				WithStartupTimeout(10*time.Second),
		),
	)

	if err != nil {
		t.Fatalf("failed to start postgres container: %s", err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx)

	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	return connStr
}
