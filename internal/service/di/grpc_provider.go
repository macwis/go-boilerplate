package di

import (
	"net"

	"github.com/google/wire"

	"github.com/macwis/go-boilerplate/internal/service/config"
)

var GRPCProvider = wire.NewSet( //nolint:gochecknoglobals
	newGRPCNetListner,
)

func newGRPCNetListner(
	cfg *config.Config,
) (net.Listener, error) {
	return net.Listen("tcp", ":"+cfg.GRPCPort)
}
