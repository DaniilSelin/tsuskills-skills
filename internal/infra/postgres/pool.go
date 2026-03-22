package postgres

import (
	"context"
	"fmt"
	"time"

	"tsuskills-skills/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, cfg *config.PostgresConfig) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("%s&search_path=%s,public", cfg.Pool.ConnConfig.ConnString(), cfg.Schema)
	poolCfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parse pool config: %w", err)
	}
	poolCfg.MaxConns = cfg.Pool.MaxConns
	poolCfg.MinConns = cfg.Pool.MinConns
	poolCfg.MaxConnLifetime = cfg.Pool.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.Pool.MaxConnIdleTime
	poolCfg.HealthCheckPeriod = cfg.Pool.HealthCheckPeriod

	for attempt := 0; attempt <= cfg.ConnectRetries; attempt++ {
		pool, connErr := pgxpool.NewWithConfig(ctx, poolCfg)
		if connErr == nil {
			if pool.Ping(ctx) == nil {
				return pool, nil
			}
			pool.Close()
		}
		if attempt == cfg.ConnectRetries {
			break
		}
		select {
		case <-time.After(cfg.ConnectRetryDelay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return nil, fmt.Errorf("failed to connect after %d attempts", cfg.ConnectRetries+1)
}

func RunMigrations(databaseURL, migrationsPath string) error {
	return runMigrate(databaseURL, migrationsPath)
}
