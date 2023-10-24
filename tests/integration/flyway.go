package integration

import (
	"context"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	FlywayContainerName = "myapp-test-migrations"
	FlywayImage         = "flyway/flyway:latest"
	MigrationsName      = "test-migrations"
)

func NewFlywayMigrationService(
	t *testing.T,
	ctx context.Context,
	dockerNetwork *testcontainers.DockerNetwork,
) testcontainers.Container {
	reqFlyway := testcontainers.ContainerRequest{
		Name:       lo.Ternary(ReuseContainers, FlywayContainerName, ""),
		Image:      FlywayImage,
		AutoRemove: false,
		WaitingFor: wait.ForAll(
			wait.ForLog("Successfully applied"),
			wait.ForExit(),
		),
		Mounts: testcontainers.Mounts(
			testcontainers.BindMount(ProjectRoot+"/sql/migrations/", "/flyway/sql/"),
			testcontainers.BindMount(ProjectRoot+"/sql/migrations/flyway.conf", "/flyway/conf/flyway.conf"),
		),
		Env: map[string]string{
			"POSTGRES_USER":           PgUser,
			"POSTGRES_PASSWORD":       PgPassword,
			"POSTGRES_DB":             PgDB,
			"APPLICATION_DB_HOST":     PgHost,
			"APPLICATION_DB":          PgAppDB,
			"APPLICATION_DB_USERNAME": PgAppUser,
			"APPLICATION_DB_PASSWORD": PgAppPassword,
		},
		Cmd:            []string{"migrate"},
		Networks:       []string{dockerNetwork.Name},
		NetworkAliases: map[string][]string{dockerNetwork.Name: {MigrationsName}},
	}
	flyway, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: reqFlyway,
		Started:          true,
	})
	require.NoError(t, err)

	return flyway
}
