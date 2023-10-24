package di

import (
	"github.com/google/wire"

	"github.com/macwis/go-boilerplate/internal/service/config"
)

var CommonProvider = wire.NewSet( //nolint:gochecknoglobals\
	config.NewLogger,
	config.NewCancelChannel,
	config.NewContext,
)
