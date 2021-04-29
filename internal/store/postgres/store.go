package postgres

import (
	"context"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver in use.
	"github.com/nairobi-gophers/fupisha/internal/store"
	"github.com/pkg/errors"
)

//Store is a postgresql implementation of our store interface
type Store struct {
	db        *sqlx.DB
	userStore *userStore
	urlStore  *urlStore
}

//Users returns a user store.
func (s *Store) Users() store.UserStore {
	return s.userStore
}

//Urls returns a url store.
func (s *Store) Urls() store.URLStore {
	return s.urlStore
}

//Migrate migrates the store database schema.
func (s *Store) Migrate() error {
	for _, q := range migrate {
		_, err := s.db.Exec(q)
		if err != nil {
			return errors.Wrap(err, "migrating schema")
		}
	}
	return nil
}

//Drop drops the store database schema.
func (s *Store) Drop() error {
	for _, q := range drop {
		_, err := s.db.Exec(q)
		if err != nil {
			return errors.Wrap(err, "dropping schema")
		}
	}
	return nil
}

// Reset resets the store database to its initial state.
func (s *Store) Reset() error {
	err := s.Drop()
	if err != nil {
		return err
	}
	return s.Migrate()
}

var _ store.Store = (*Store)(nil) //Validate that store object actually points to something.

//Connect connects to a postgres store and returns an initialized postgres store object.
//address: localhost:5432
func Connect(address, username, password, database string) (*Store, error) {
	sslMode := "disable" //Should be set in the config object
	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")
	q.Set("connect_timeout", "10")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(username, password),
		Host:     address,
		Path:     database,
		RawQuery: q.Encode(),
	}

	db, err := sqlx.Open("postgres", u.String())
	if err != nil {
		return nil, errors.Wrap(err, "connecting to database")
	}

	if err := statusCheck(context.Background(), db); err != nil {
		return nil, errors.Wrap(err, "connect: connection never ready")
	}

	s := Store{
		db:        db,
		userStore: &userStore{db: db},
		urlStore:  &urlStore{db: db},
	}

	err = s.Migrate()
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func statusCheck(ctx context.Context, db *sqlx.DB) error {
	var pingError error
	for attempts := 1; ; attempts++ {
		pingError = db.Ping()
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	//I am  paranoid and we like to detect any false positive.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}
