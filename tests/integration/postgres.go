package integration

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/uptrace/bun/driver/pgdriver"
)

type PostgresServer struct {
	instance         testcontainers.Container
	migrationService testcontainers.Container
}

func (p PostgresServer) Port(t *testing.T) int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	port, err := p.instance.MappedPort(ctx, PgPort)
	require.NoError(t, err)
	return port.Int()
}

func (p PostgresServer) Close(t *testing.T) {
	if !ReuseContainers {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		require.NoError(t, p.instance.Terminate(ctx))
	}
}

func (p PostgresServer) Address(t *testing.T) string {
	return fmt.Sprintf("postgres://%s:%s@127.0.0.1:%d/%s?sslmode=disable",
		PgAppUser, PgAppPassword, p.Port(t), PgAppDB)
}

func (p PostgresServer) Prep(t *testing.T) {
	dsn := p.Address(t)
	t.Logf("postgres: %s", dsn)
	sqlConn := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	err := sqlConn.PingContext(context.TODO())
	require.NoError(t, err)
}

const (
	PgContainerName = "myapp-postgres"
	PgImage         = "postgres:10"
	PgHost          = "myapp-db"
	PgUser          = "postgres"
	PgPassword      = "postgres"
	PgDB            = "postgres"
	PgPort          = "5432"
	PgAppUser       = "myapp_test_user"
	PgAppPassword   = "test_user_password"
	PgAppDB         = "myapp_test_db"
)

func NewPostgresServer(t *testing.T, ctx context.Context, dockerNetwork *testcontainers.DockerNetwork) *PostgresServer {
	t.Helper()

	ctx, cancel := context.WithTimeout(ctx, ContainerContextTimeout)
	defer cancel()

	if dockerNetwork == nil {
		dockerNetwork = CreateDockerNetwork(t, ctx, nil)
	}

	defer dockerNetwork.Remove(ctx) //nolint:errcheck

	port := fmt.Sprintf("%s/tcp", PgPort)
	req := testcontainers.ContainerRequest{
		Name:         lo.Ternary(ReuseContainers, PgContainerName, ""),
		Hostname:     PgHost,
		Image:        PgImage,
		ExposedPorts: []string{port},
		AutoRemove:   false,
		SkipReaper:   ReuseContainers,
		Env: map[string]string{
			"POSTGRES_USER":           PgUser,
			"POSTGRES_PASSWORD":       PgPassword,
			"POSTGRES_DB":             PgDB,
			"APPLICATION_DB":          PgAppDB,
			"APPLICATION_DB_USERNAME": PgAppUser,
			"APPLICATION_DB_PASSWORD": PgAppPassword,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort(nat.Port(port)),
		),
		Mounts: testcontainers.Mounts(
			testcontainers.BindMount(ProjectRoot+"/sql/init", "/docker-entrypoint-initdb.d"),
		),
		Networks:       []string{dockerNetwork.Name},
		NetworkAliases: map[string][]string{dockerNetwork.Name: {PgHost}},
		Cmd:            []string{"postgres", "-c", "log_statement=all", "-c", "log_destination=stderr"},
	}
	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            ReuseContainers,
	})
	require.NoError(t, err)

	flyway := NewFlywayMigrationService(t, ctx, dockerNetwork)

	return &PostgresServer{
		instance:         postgres,
		migrationService: flyway,
	}
}
