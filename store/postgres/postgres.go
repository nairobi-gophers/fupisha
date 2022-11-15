package postgres

import (
	"context"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver in use.
	"github.com/nairobi-gophers/fupisha/store"
	"github.com/pkg/errors"
)

// compilation check for store.Store concrete implementation.
var _ store.Store = (*Store)(nil)

// NewStore creates and returns an initialized postgresql store for use as our state backend.
func NewStore(cfg *Config) (*Store, error) {
	db, err := connect(cfg)
	if err != nil {
		return nil, err
	}

	s := Store{
		&userStore{db: db},
		&urlStore{db: db},
	}

	err = migrateState(db)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// connects to a postgres store and returns an initialized postgres store object.
// address: localhost:5432
func connect(cfg *Config) (*sqlx.DB, error) {
	sslMode := "disable" //Should be set in the config object
	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")
	q.Set("connect_timeout", "10")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	db, err := sqlx.Open("postgres", u.String())
	if err != nil {
		return nil, errors.Wrap(err, "connecting to database")
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()

	if err := statusCheck(ctx, db); err != nil {
		return nil, errors.Wrap(err, "connect: connection never ready")
	}

	return db, nil
}

type Config struct {
	//Host e.g. localhost:5432
	Host string
	//Password is the database user's password
	Password string
	//User is the database username
	User string
	//Name is the database name
	Name string
	//MaxIdleConns is the maximum number of conns in the idle conn pool
	MaxIdleConns int
	//MaxOpenConns is the maximum number of open conns to the database.
	MaxOpenConns int
	//DisableTLS enable TLS on connections to the database.
	DisableTLS bool
}

// Store is a postgresql implementation of our store interface
type Store struct {
	*userStore
	*urlStore
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

// migrates the store database schema.
func migrateState(db *sqlx.DB) error {
	for _, q := range migrate {
		_, err := db.Exec(q)
		if err != nil {
			return errors.Wrap(err, "migrating schema")
		}
	}
	return nil
}

// drops the store database schema.
func dropState(db *sqlx.DB) error {
	for _, q := range drop {
		_, err := db.Exec(q)
		if err != nil {
			return errors.Wrap(err, "dropping schema")
		}
	}
	return nil
}

// resets the store database to its initial state.
//func resetState(db *sqlx.DB) error {
//	err := dropState(db)
//	if err != nil {
//		return err
//	}
//	return migrateState(db)
//}
