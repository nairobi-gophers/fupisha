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
		t.Fatal(err)
	}

	err = s.Reset()
	if err != nil {
		t.Fatal(err)
	}

	//tear down a table.
	tearDown := func() {
		err := s.Drop()
		if err != nil {
			t.Fatal(err)
		}
	}

	return s, tearDown
}
