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

	testContainer.Create(t)

	purgeContainer := func() {
		t.Logf("Purging test container...")
		//purge the test container
		if err := testContainer.resource.Close(); err != nil {
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

	db, err := connect(opts)
	if err != nil {
		t.Fatal(err)
	}

	if err := statusCheck(ctx, db); err != nil {
		t.Fatalf("status check database: %s", err)
	}

	err = migrateState(db)
	if err != nil {
		t.Fatal(err)
	}

	closedb := func() {
		t.Log("Closing database connection...")
		//close database connection
		db.Close()
	}

	dropdb := func() {
		t.Log("Dropping database...")
		//drop the database
		err := dropState(db)
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

	return &Store{
		&userStore{db: db},
		&urlStore{db: db},
	}, teardown
}
