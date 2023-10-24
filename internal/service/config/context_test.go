package config

import (
	"syscall"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestNewContext(t *testing.T) {
	t.Parallel()

	cancelChannel := NewCancelChannel()
	ctx := NewContext(logrus.New(), cancelChannel)

	cancelChannel <- struct{}{}
	<-ctx.Done()
}

func TestQuitSignal(t *testing.T) {
	t.Parallel()

	cancelChannel := NewCancelChannel()
	_ = NewContext(logrus.New(), cancelChannel)

	err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	require.NoError(t, err)
}
