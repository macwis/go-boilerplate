package e2e

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/macwis/go-boilerplate/tests/integration"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"testing"
	"time"

	docker "github.com/docker/docker/client"
)

func PrivateKeyHexToAccountAddress(privateKeyHex string) (common.Address, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return common.Address{}, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return address, nil
}

func PauseContainer(container testcontainers.Container, delay time.Duration, t *testing.T) {
	dockerClient, err := docker.NewClientWithOpts(docker.FromEnv)
	require.NoError(t, err)

	tag := container.GetContainerID()
	err = dockerClient.ContainerPause(context.TODO(), tag)
	require.NoError(t, err)

	time.Sleep(delay)

	err = dockerClient.ContainerUnpause(context.TODO(), tag)
	require.NoError(t, err)
}

func GetHTTPUrl(t *testing.T, ganache *integration.GanacheServer) (string, error) {
	if ganache != nil {
		return fmt.Sprintf("ws://%s", ganache.Address(t)), nil
	}
	return "", fmt.Errorf("instance nil")
}
