package e2e

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"testing"
)

func (s *MyAppTestSuite) TestExample() {
	log := s.service.AppInstance.GetLogger()
	cfg := s.service.TestCluster.Cfg
	defer s.service.CancelFunc()
	stopApp := s.startApplication()
	defer stopApp()

	fromAddress, err := s.service.TestCluster.GetAccount(0)
	require.NoError(s.T(), err)

	toAddress, err := s.service.TestCluster.GetPrivateKey(0)
	require.NoError(s.T(), err)

	client, err := ethclient.Dial(GetHTTPUrl(s.T(), s.service.TestCluster.TestCluster.Ganache.Address()))

	s.T().Run("send transaction", func(t *testing.T)) {
		require.Equal(s.T(), 0, CheckBalance())
	}
}
