package postgres

import (
	"context"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/pkg/errors"
)

func newTestDatabase(tb testing.TB) (*Store, error) {
	tb.Helper()

	ctx := context.Background()

	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, err
	}

	testContainer := NewPostgresqlContainer(pool)

	if err := testContainer.Create(); err != nil {
		return nil, err
	}

	db := testContainer.Connect()

	tb.Cleanup(func() {

		//close connection
		db.Close()
	})

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
			tb.Fatal(err)
		}
	}

	tb.Cleanup(tearDown)

	return s, nil
}