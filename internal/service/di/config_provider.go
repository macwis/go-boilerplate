package di

import (
	"github.com/google/wire"

	"github.com/macwis/go-boilerplate/internal/service/config"
)

var ConfigProvider = wire.NewSet( //nolint:gochecknoglobals
	config.NewConfig,
)
