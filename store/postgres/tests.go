package postgres

import (
	"context"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
)

func NewTestDatabase(t *testing.T) (*Store, func()) {

	ctx := context.Background()

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatal(err)
	}

	testContainer := NewPostgresqlContainer(pool)

	resource, err := testContainer.Create()
	if err != nil {
		t.Fatal(err)
	}

	testContainer.resource = resource

	purgeContainer := func() {
		//purge the test container
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("Could not purge resource: %s", err)
		}
	}
	opts := &Config{
		Host:       "localhost:5432",
		User:       "testcontainer",
		Password:   "Aa123456.",
		Name:       "testcontainer",
		DisableTLS: true,
	}

	s, err := Connect(opts)
	if err != nil {
		t.Fatal(err)
	}

	if err := statusCheck(ctx, s.db); err != nil {
		t.Fatalf("status check database: %s", err)
	}

	closedb := func() {
		//close database connection
		s.db.Close()
	}

	dropdb := func() {
		//drop the database
		err := s.Drop()
		if err != nil {
			t.Fatal(err)
		}
	}
	//tear down a table.
	teardown := func() {
		t.Helper()
		dropdb()
		closedb()
		purgeContainer()
	}

	return s, teardown
}
