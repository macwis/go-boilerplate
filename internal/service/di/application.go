package di

import (
	"context"
	"database/sql"
	"net"
	"os"
	"os/signal"

	"github.com/macwis/go-boilerplate/internal/service/config"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type WaitGroup interface {
	Add(int)
	Done()
	Wait()
}

type Application struct {
	ctx    context.Context
	log    *logrus.Logger
	config *config.Config

	db          *sql.DB
	netListener net.Listener
	// grpcServer  *grpc.Server
}

func (app *Application) Run() error {
	var _, cancel = context.WithCancel(app.ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	defer func() {
		signal.Stop(c)
		cancel()
	}()

	errGroup, _ := errgroup.WithContext(app.ctx)

	//errGroup.Go(func() error {
	//	// TODO: RegisterServer()
	//	return app.grpcServer.Serve(app.netListener)
	//})

	err := errGroup.Wait()

	app.log.Info("app stopped")

	return err
}

func (app *Application) ShutdownAndCleanup() {
	app.log.Info("app shutting down")

	// TODO: all graceful shutdown actions
}

func (app *Application) GetLogger() *logrus.Logger {
	return app.log
}

func (app *Application) GetContext() context.Context {
	return app.ctx
}

func (app *Application) GetConfig() *config.Config {
	return app.config
}

type TestApplication struct {
	*Application

	Ctx    context.Context
	Log    *logrus.Logger
	Config *config.Config

	DB          *sql.DB
	netListener net.Listener
	// GRPCServer  *grpc.Server
}
