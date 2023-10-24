//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"

	"github.com/macwis/go-boilerplate/internal/service/config"
)

func SetupApplication() (*Application, func(), error) {
	panic(wire.Build(wire.NewSet(
		CommonProvider,
		ConfigProvider,
		GRPCProvider,
		StorageProvider,

		wire.Struct(new(Application), "*"))),
	)
}

func SetupApplicationForIntegrationTests(cfg *config.Config) (*TestApplication, func(), error) {
	panic(wire.Build(wire.NewSet(
		CommonProvider,
		GRPCProvider,
		StorageProvider,

		wire.Struct(new(Application), "*"),
		wire.Struct(new(TestApplication), "*"))),
	)
}
