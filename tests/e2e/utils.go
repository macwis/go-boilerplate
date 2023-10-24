package e2e

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"

	"github.com/macwis/go-boilerplate/tests/integration"

	docker "github.com/docker/docker/client"
)

const (
	TestGasLimit             = 21000      // 21 Kwei
	TestMaxPriorityFeePerGas = 2000000000 // 2 Gwei
	TestMaxFeePerGas         = 2000000000 // 20 Gwei
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

func GetHTTPUrl(t *testing.T, ganache integration.ServerInstance) (string, error) {
	if ganache != nil {
		return fmt.Sprintf("ws://%s", ganache.Address(t)), nil
	}
	return "", fmt.Errorf("instance nil")
}

func CheckBalance(client *ethclient.Client, address common.Address) (*big.Int, error) {
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func NewTransaction(
	client *ethclient.Client,
	privateKeyHex string,
	value *big.Int,
	toAddress common.Address,
) (*types.Transaction, error) {
	fromAddress, err := PrivateKeyHexToAccountAddress(privateKeyHex)
	if err != nil {
		return nil, err
	}

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	gasLimit := uint64(TestGasLimit)
	tip := big.NewInt(TestMaxPriorityFeePerGas)
	feeCap := big.NewInt(TestMaxFeePerGas)

	var data []byte

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: tip,
		GasFeeCap: feeCap,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     value,
		Data:      data,
	})

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func SendTransaction(client *ethclient.Client, signedTx *types.Transaction) (*types.Transaction, error) {
	err := client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}
