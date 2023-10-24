package integration

import (
	"context"
	"fmt"
	"github.com/macwis/go-boilerplate/internal/myapp/utils"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	GanacheContainerName = "myapp-ganache"
	GanacheImage         = "trufflesuite/ganache"
	GanacheHost          = "ganache"
	GanachePort          = "8545"
)

type GanacheServer struct {
	Instance testcontainers.Container
}

func (g GanacheServer) Port(t *testing.T) int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	port, err := g.Instance.MappedPort(ctx, GanachePort)
	require.NoError(t, err)
	return port.Int()
}

func (g GanacheServer) Close(t *testing.T) {
	if !ReuseContainers {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		require.NoError(t, g.Instance.Terminate(ctx))
	}
}

func (g GanacheServer) Address(t *testing.T) string {
	return fmt.Sprintf("127.0.0.1:%d", g.Port(t))
}

func (g GanacheServer) Prep(t *testing.T) {
	t.Logf("ganache: %s", g.Address(t))
}

func NewGanacheServer(
	t *testing.T,
	ctx context.Context,
	dockerNetwork *testcontainers.DockerNetwork,
) *GanacheServer {
	t.Helper()

	ctx, cancel := context.WithTimeout(ctx, ContainerContextTimeout)
	defer cancel()

	port := fmt.Sprintf("%s/tcp", GanachePort)
	request := testcontainers.ContainerRequest{
		Name:           lo.Ternary(ReuseContainers, GanacheContainerName, ""),
		Hostname:       GanacheHost,
		Image:          GanacheImage,
		ExposedPorts:   []string{port},
		AutoRemove:     false,
		SkipReaper:     ReuseContainers,
		Env:            map[string]string{},
		Networks:       []string{dockerNetwork.Name},
		NetworkAliases: map[string][]string{dockerNetwork.Name: {GanacheHost}},
		Entrypoint: []string{
			"node",
			"/app/dist/node/cli.js",
			"--server.host", "0.0.0.0",
			"--chain.chainId=1337",
			"--chain.networkId=1337",
			"--miner.blockTime=3",
			"--logging.debug",
			"--logging.verbose",
			"--wallet.accounts",
			fmt.Sprintf("0x12aaa41f8c755a8576fb67122e4d167bf10083f9381372755c522d778f4993e3,%s", utils.Eth_1K),
			fmt.Sprintf("0xf893af1ff2cc1a46915cf9ff7d038f5a6b3472c9bf92ac43f75088479a192c6f,%s", utils.Eth_1K),
			fmt.Sprintf("0xf893af1ff2cc1a46915cf9ff7d038f5a6b3472c9bf92ac43f75088479a192c6f,%s", utils.Eth_1K),
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("RPC Listening on 0.0.0.0:8545"),
			wait.ForListeningPort(nat.Port(port)),
		),
	}
	instance, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
		Reuse:            ReuseContainers,
	})
	require.NoError(t, err)

	return &GanacheServer{
		Instance: instance,
	}
}
