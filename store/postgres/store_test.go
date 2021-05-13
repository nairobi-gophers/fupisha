package postgres

import (
	"context"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/pkg/errors"
)

func newTestDatabase(t *testing.T) (*Store, error) {
	t.Helper()

	ctx := context.Background()

	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, err
	}

	testContainer := NewPostgresqlContainer(pool)

	resource, err := testContainer.Create()
	if err != nil {
		return nil, err
	}

	testContainer.resource = resource

	purgeContainer := func() {
		//purge the test container
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("Could not purge resource: %s", err)
		}
	}

	t.Cleanup(purgeContainer)

	db := testContainer.Connect()

	closeDB := func() {
		//close connection
		db.Close()
	}

	t.Cleanup(closeDB)

	if err := statusCheck(ctx, db); err != nil {
		return nil, errors.Wrap(err, "status check database: %s")
	}

	s := &Store{
		db:        db,
		userStore: &userStore{db: db},
		urlStore:  &urlStore{db: db},
	}

	if err := s.Migrate(); err != nil {
		return nil, errors.Wrap(err, "failed to migrate database")
	}

	//tear down a table.
	tearDown := func() {
		err := s.Drop()
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Cleanup(tearDown)

	return s, nil
}
