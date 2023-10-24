package e2e

import (
	"github.com/macwis/go-boilerplate/internal/myapp/utils"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func (s *MyAppTestSuite) TestExample() {
	// log := s.service.AppInstance.GetLogger()
	cfg := s.service.TestCluster.Cfg
	var err error
	defer s.service.CancelFunc()
	stopApp := s.startApplication()
	defer stopApp()

	fromAddress, err := s.service.TestCluster.GetAccount(0)
	fromPrvKey, err := s.service.TestCluster.GetPrivateKey(0)
	require.NoError(s.T(), err)

	toAddress, err := s.service.TestCluster.GetAccount(1)
	require.NoError(s.T(), err)

	client, err := ethclient.Dial(cfg.GethURL)
	require.NoError(s.T(), err)

	s.T().Run("send transaction", func(t *testing.T) {
		fromBalance, err := CheckBalance(client, fromAddress)
		require.NoError(s.T(), err)

		toBalance, err := CheckBalance(client, toAddress)
		require.NoError(s.T(), err)

		require.Equal(s.T(), 0, fromBalance.Cmp(utils.Eth_1K), "from address incorrect initial balance")
		require.Equal(s.T(), 0, toBalance.Cmp(utils.Eth_1K), "from address incorrect initial balance")

		txSigned, err := NewTransaction(client, fromPrvKey, utils.Mwei, toAddress)
		require.NoError(t, err)

		_, err = SendTransaction(client, txSigned)
		require.NoError(t, err)

		fromBalance2, err := CheckBalance(client, fromAddress)
		require.NoError(s.T(), err)

		// TODO: manually call mine or make sure it's mined

		require.Equal(s.T(), 0, fromBalance2.Cmp(utils.Eth_1K), "from address incorrect initial balance")
	})
}
