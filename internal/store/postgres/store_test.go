package postgres

import (
	"testing"
)

const (
	testUsername = "fp_user"
	testPassword = "fp_s35r37_!"
	testDatabase = "fupisha"
	testAddress  = "localhost:5432"
)

func testConn(t *testing.T) (*Store, func()) {
	s, err := Connect(testAddress, testUsername, testPassword, testDatabase)
	if err != nil {
		t.Fatalf("failed to connect to the test postgres database: address=%q, username=%q, password=%q, database=%q: %s", testAddress, testUsername, testPassword, testDatabase, err)
	}

	err = s.Reset()
	if err != nil {
		t.Fatalf("failed to reset the test postgresql database: %s", err)
	}

	//tear down a table.
	tearDown := func() {
		err := s.Drop()
		if err != nil {
			t.Fatalf("failed to drop table of the test postgres database: %s", err)
		}
	}

	return s, tearDown
}
