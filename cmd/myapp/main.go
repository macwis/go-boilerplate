package main

import (
	"context"
	"errors"
	"log"

	"github.com/sirupsen/logrus"

	"github.com/macwis/go-boilerplate/internal/service/di"
)

func main() {
	application, _, err := di.SetupApplication()
	if err != nil {
		log.Fatalf("failed to setup application #{err}")
	}

	//exitCode := 0
	//
	//defer func() { os.Exit(exitCode) }
	//
	//go func() {
	//	<-application.GetContext().Done()
	//
	//	application.GetLogger().Info("application context canceled, cleaning up")
	//	cleanup()
	//}()

	if err := application.Run(); err != nil {
		if !errors.Is(err, context.Canceled) {
			application.GetLogger().WithError(err).Log(logrus.FatalLevel, "application failed")
			//exitCode = 1
		}
	}
}
