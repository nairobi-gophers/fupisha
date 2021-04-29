package postgres

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

type testConfig struct {
	Host     string
	User     string
	Password string
	Database string
	Port     int
}

type PostgresqlContainer struct {
	pool     *dockertest.Pool
	resource *dockertest.Resource
	image    string
	opts     testConfig
}

func NewPostgresqlContainer(pool *dockertest.Pool) PostgresqlContainer {
	opts := testConfig{
		Host:     "localhost",
		User:     "testcontainer",
		Password: "Aa123456.",
		Database: "testcontainer",
		Port:     5432,
	}

	return PostgresqlContainer{pool: pool, opts: opts, image: "postgresql-testcontainer"}
}

func (container PostgresqlContainer) Create() error {
	if isRunning(*container.pool, container.image) {
		return nil
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
		Name: container.image,
	}

	resource, err := container.pool.RunWithOptions(&dockerOpts)
	if err != nil {
		log.Fatalf("Could not start resource (Postgresql Test Container): %s", err.Error())
		return err
	}

	container.resource = resource
	return nil
}

func (container PostgresqlContainer) Connect() *sqlx.DB {
	var db *sqlx.DB
	if err := container.pool.Retry(func() error {
		defaultDsn := "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"
		dsn := fmt.Sprintf(defaultDsn, container.opts.Host, container.opts.User, container.opts.Password, container.opts.Database, container.opts.Port)

		var err error
		db, err = sqlx.Open("postgres", dsn)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return db
}

func isRunning(pool dockertest.Pool, imagename string) bool {
	dockerContainers, _ := pool.Client.ListContainers(docker.ListContainersOptions{
		All: false,
	})

	for _, dockerContainer := range dockerContainers {
		for _, name := range dockerContainer.Names {
			if strings.Contains(name, imagename) {
				// fmt.Printf("%s image is running...", dockerContainer.Image)
				return true
			}
		}
	}

	return false
}
