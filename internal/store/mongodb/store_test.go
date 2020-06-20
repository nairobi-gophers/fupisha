package mongodb

import (
	"context"
	"testing"
)

const (
	testUsername = "fupisha"
	testPassword = "fupisha"
	testDatabase = "fupisha"
	testAddress  = "localhost:27017"
)

func testConn(t *testing.T) (*Store, func(col string)) {
	s, err := Connect(testAddress, testUsername, testPassword, testDatabase)
	if err != nil {
		t.Fatalf(
			"failed to connect to the test mongo database: address=%q, username=%q, password=%q, database=%q: %s",
			testAddress, testUsername, testPassword, testDatabase, err,
		)
	}

	//tear down a collection.
	tearDown := func(col string) {
		err := s.db.Collection(col).Drop(context.Background())
		if err != nil {
			t.Fatalf("failed to drop collection of the test mongo database: %s", err)
		}
	}

	return s, tearDown
}
