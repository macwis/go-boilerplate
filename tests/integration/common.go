package integration

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go"

	"github.com/samber/lo"
)

var (
	_, _file, _, _  = runtime.Caller(0)                           // nolint:gochecknoglobals
	ProjectRoot     = filepath.Join(filepath.Dir(_file), "../..") // nolint:gochecknoglobals
	ReuseContainers = os.Getenv("REUSE_CONTAINERS") == "1"        // nolint:gochecknoglobals
)

const (
	ContainerContextTimeout = 1 * time.Minute
)

func CreateDockerNetwork(t *testing.T, ctx context.Context, name *string) *testcontainers.DockerNetwork {
	if name == nil {
		if ReuseContainers {
			name = lo.ToPtr("myapp-test")
		} else {
			name = lo.ToPtr(uuid.NewString())
		}
	}
	t.Logf("docker network name: %s", lo.FromPtr(name))

	network, err := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{Name: *name},
	})

	require.NoError(t, err)

	t.Cleanup(func() {
		_ = network.Remove(context.Background())
	})

	return network.(*testcontainers.DockerNetwork)
}

type TestCluster struct {
	t   *testing.T
	ctx context.Context

	network  *testcontainers.DockerNetwork
	Postgres ServerInstance
	Redis    ServerInstance
	Ganache  ServerInstance
	OTEL     ServerInstance
	Teardown func()
}

func (tc *TestCluster) Close(t *testing.T) {
	if tc.Postgres != nil {
		tc.Postgres.Close(t)
	}

	if tc.Redis != nil {
		tc.Redis.Close(t)
	}

	if tc.Ganache != nil {
		tc.Ganache.Close(t)
	}

	if tc.OTEL != nil {
		tc.OTEL.Close(t)
	}
}

func NewTestCluster(t *testing.T, ctx context.Context, options ...func(*TestCluster)) *TestCluster {
	t.Helper()

	network := CreateDockerNetwork(t, ctx, nil)

	tc := &TestCluster{
		t:       t,
		ctx:     ctx,
		network: network,
	}

	for _, o := range options {
		o(tc)
	}

	return tc
}

func WithAll() func(cluster *TestCluster) {
	return func(tc *TestCluster) {
		tc.Postgres = NewPostgresServer(tc.t, tc.ctx, tc.network)
		tc.Ganache = NewGanacheServer(tc.t, tc.ctx, tc.network)
		//tc.Redis = NewRedisServer(tc.t, tc.ctx, tc.network)
		//tc.OTEL = NewOTELServer(tc.t, tc.ctx, tc.network)

		tc.Postgres.Prep(tc.t)
		tc.Ganache.Prep(tc.t)
		//tc.Redis.Prep(tc.t)
		//tc.OTEL.Prep(tc.t)
	}
}

func WithPostgres() func(*TestCluster) {
	return func(tc *TestCluster) {
		tc.Postgres = NewPostgresServer(tc.t, tc.ctx, tc.network)
		tc.Postgres.Prep(tc.t)
	}
}

func WithGanache() func(*TestCluster) {
	return func(tc *TestCluster) {
		tc.Ganache = NewGanacheServer(tc.t, tc.ctx, tc.network)
		tc.Ganache.Prep(tc.t)
	}
}
