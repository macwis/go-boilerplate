package storage

import (
	"context"
	"database/sql"
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/macwis/go-boilerplate/internal/service/config"
)

func NewSQLDatabase(ctx context.Context, cfg *config.Config, log *logrus.Logger) (*sql.DB, error) {
	_, err := url.Parse(cfg.DatastoreURL)
	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.DatastoreURL)))

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
