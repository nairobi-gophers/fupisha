package postgres

import (
	"strconv"
	"strings"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"
)

type testConfig struct {
	Host     string
	User     string
	Password string
	Database string
	Port     int
}

type PostgresqlContainer struct {
	pool      *dockertest.Pool
	resource  *dockertest.Resource
	imagename string
	opts      testConfig
}

func NewPostgresqlContainer(pool *dockertest.Pool) *PostgresqlContainer {
	opts := testConfig{
		Host:     "localhost",
		User:     "testcontainer",
		Password: "Aa123456.",
		Database: "testcontainer",
		Port:     5432,
	}

	return &PostgresqlContainer{pool: pool, opts: opts, imagename: "postgresql-testcontainer"}
}

func (container *PostgresqlContainer) Create(t *testing.T) {

	if isRunning(container.pool, container.imagename) {
		t.Fatal(errors.New("container already exists and is running"))
	}

	dockerOpts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_USER=" + container.opts.User,
			"POSTGRES_PASSWORD=" + container.opts.Password,
			"POSTGRES_DB=" + container.opts.Database,
		},
		ExposedPorts: []string{strconv.Itoa(container.opts.Port)},
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port(strconv.Itoa(container.opts.Port)): {{HostIP: "0.0.0.0", HostPort: strconv.Itoa(container.opts.Port)}},
		},
		Name: container.imagename,
	}

	resource, err := container.pool.RunWithOptions(&dockerOpts, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		t.Fatal(err)
	}

	container.resource = resource
}

func isRunning(pool *dockertest.Pool, imagename string) bool {
	dockerContainers, _ := pool.Client.ListContainers(docker.ListContainersOptions{
		All: true,
	})

	for _, dockerContainer := range dockerContainers {
		for _, name := range dockerContainer.Names {
			if strings.Contains(name, imagename) {
				return true
			}
		}
	}
	return false
}
