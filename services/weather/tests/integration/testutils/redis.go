package testutils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestRedis struct {
	container testcontainers.Container
	client    *redis.Client
	host      string
	port      string
	password  string
	db        int
}

func SetupTestRedis(t *testing.T) *TestRedis {
	t.Helper()

	const password = ""
	const db = 0

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
		Cmd:          []string{"redis-server", "--appendonly", "yes"},
	}

	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := redisContainer.Host(ctx)
	require.NoError(t, err)

	mappedPort, err := redisContainer.MappedPort(ctx, "6379")
	require.NoError(t, err)

	port := mappedPort.Port()

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Ping(ctx).Result()
	require.NoError(t, err)

	testRedis := &TestRedis{
		container: redisContainer,
		client:    client,
		host:      host,
		port:      port,
	}

	t.Cleanup(func() {
		err := client.Close()
		if err != nil {
			t.Errorf("Failed to close connection with redis: %s", err.Error())
		}
		err = redisContainer.Terminate(context.Background())
		if err != nil {
			t.Errorf("Failed to terminate redis container: %s", err.Error())
		}
	})

	return testRedis
}

func (tr *TestRedis) ConnectionString() string {
	return fmt.Sprintf("redis://:%s@%s:%s/%d", tr.password, tr.host, tr.port, tr.db)
}

func (tr *TestRedis) Size() int {
	result := tr.client.DBSize(context.Background())
	return int(result.Val())
}
