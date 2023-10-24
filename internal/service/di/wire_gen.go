// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/macwis/go-boilerplate/internal/myapp/infra/storage"
	"github.com/macwis/go-boilerplate/internal/service/config"
)

// Injectors from inject.go:

func SetupApplication() (*Application, func(), error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, nil, err
	}
	logger := config.NewLogger(configConfig)
	cancelChannel := config.NewCancelChannel()
	context := config.NewContext(logger, cancelChannel)
	db, err := storage.NewSQLDatabase(context, configConfig, logger)
	if err != nil {
		return nil, nil, err
	}
	listener, err := newGRPCNetListner(configConfig)
	if err != nil {
		return nil, nil, err
	}
	application := &Application{
		ctx:         context,
		log:         logger,
		config:      configConfig,
		db:          db,
		netListener: listener,
	}
	return application, func() {
	}, nil
}

func SetupApplicationForIntegrationTests(cfg *config.Config) (*TestApplication, func(), error) {
	logger := config.NewLogger(cfg)
	cancelChannel := config.NewCancelChannel()
	context := config.NewContext(logger, cancelChannel)
	db, err := storage.NewSQLDatabase(context, cfg, logger)
	if err != nil {
		return nil, nil, err
	}
	listener, err := newGRPCNetListner(cfg)
	if err != nil {
		return nil, nil, err
	}
	application := &Application{
		ctx:         context,
		log:         logger,
		config:      cfg,
		db:          db,
		netListener: listener,
	}
	testApplication := &TestApplication{
		Application: application,
		Ctx:         context,
		Log:         logger,
		Config:      cfg,
		DB:          db,
		netListener: listener,
	}
	return testApplication, func() {
	}, nil
}