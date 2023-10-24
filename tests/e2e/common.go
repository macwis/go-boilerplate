package e2e

import (
	"context"
	"fmt"
	"github.com/macwis/go-boilerplate/internal/service/config"
	"github.com/macwis/go-boilerplate/tests/integration"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

type IntegrationTestSetup struct {
	Cfg         *config.Config
	privateKeys []string
	accounts    []common.Address
	TestCluster *integration.TestCluster
}

func (i *IntegrationTestSetup) Terminate(t *testing.T) {
	i.TestCluster.Close(t)
}

func (i *IntegrationTestSetup) GetAccount(idx int) (common.Address, error) {
	if idx >= 0 && idx < len(i.privateKeys) {
		return i.accounts[idx], nil
	}
	return common.Address{}, fmt.Errorf("invalid index %d", idx)
}

func (i *IntegrationTestSetup) GetPrivateKey(idx int) (string, error) {
	if idx >= 0 && idx < len(i.privateKeys) {
		return i.privateKeys[idx], nil
	}
	return "", fmt.Errorf("invalid index %d", idx)
}

func NewIntegrationTestSetup(t *testing.T, ctx context.Context, options func(cluster *integration.TestCluster)) *IntegrationTestSetup {
	t.Helper()

	tc := integration.NewTestCluster(t, ctx, options)

	cfg := &config.Config{
		ServiceName:    "myapp",
		ServiceVersion: "0.1",
		HTTPPort:       "",
		GRPCPort:       "",
		HealthPort:     "",
		LogLevel:       "debug",
	}

	if tc.Postgres != nil {
		cfg.DatastoreURL = tc.Postgres.Address(t)
	}

	ts := IntegrationTestSetup{
		Cfg: cfg,
		privateKeys: []string{
			"12aaa41f8c755a8576fb67122e4d167bf10083f9381372755c522d778f4993e3",
			"f893af1ff2cc1a46915cf9ff7d038f5a6b3472c9bf92ac43f75088479a192c6f",
			"0xf893af1ff2cc1a46915cf9ff7d038f5a6b3472c9bf92ac43f75088479a192c6f",
		},
		TestCluster: tc,
	}

	ts.accounts = make([]common.Address, len(ts.privateKeys))
	var err error

	for i, pk := range ts.privateKeys {
		ts.accounts[i], err = PrivateKeyHexToAccountAddress(pk)
		require.NoError(t, err)
	}

	return &ts
}
