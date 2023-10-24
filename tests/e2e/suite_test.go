package e2e

import (
	"context"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/macwis/go-boilerplate/internal/service/di"
	"github.com/macwis/go-boilerplate/tests/integration"
)

type MyAppIntegrationService struct {
	CancelFunc  context.CancelFunc
	AppInstance *di.TestApplication
	TestCluster *IntegrationTestSetup
}

type MyAppTestSuite struct {
	suite.Suite
	service MyAppIntegrationService
}

func (s *MyAppTestSuite) SetupTest() {
	s.service = initTestService(s.T())
	// TODO: s.truncateDB()
	// TODO: s.flushRedis()
}

func (s *MyAppTestSuite) TearDownTest() {
	s.service.TestCluster.Terminate(s.T())
	s.service.CancelFunc()
}

func TestMyAppServiceSuit(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(MyAppTestSuite))
}

func initTestService(t *testing.T) MyAppIntegrationService {
	ctx := context.TODO()

	testSetup := NewIntegrationTestSetup(t, ctx, integration.WithAll())

	app, cleanup, err := di.SetupApplicationForIntegrationTests(testSetup.Cfg)
	require.NoError(t, err)

	return MyAppIntegrationService{
		CancelFunc:  cleanup,
		AppInstance: app,
		TestCluster: testSetup,
	}
}

func (s *MyAppTestSuite) startApplication() func() {
	s.T().Helper()

	go func() {
		err := s.service.AppInstance.Run()
		require.NoError(s.T(), err)
	}()

	return func() {
		err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		require.NoError(s.T(), err)
	}
}
