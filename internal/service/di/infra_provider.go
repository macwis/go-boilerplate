package di

import (
	"github.com/google/wire"

	"github.com/macwis/go-boilerplate/internal/myapp/infra/storage"
)

var StorageProvider = wire.NewSet( //nolint:gochecknoglobals
	storage.NewSQLDatabase,
)
