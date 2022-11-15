package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/nairobi-gophers/fupisha/encoding"
	"github.com/nairobi-gophers/fupisha/store"
	"github.com/pkg/errors"
)

type userStore struct {
	db *sqlx.DB
}

// NewUser creates a new user record.
func (s userStore) NewUser(ctx context.Context, email, password string) (store.User, error) {

	var now time.Time = time.Now()

	var u store.User

	user := store.User{
		ID:                  encoding.GenUniqueID(),
		Email:               email,
		Password:            password,
		VerificationToken:   encoding.GenUniqueID(),
		VerificationExpires: now.Add(time.Minute * 15).UTC().Round(time.Microsecond), //expires 15 mins later
		CreatedAt:           now.UTC().Round(time.Microsecond),
		UpdatedAt:           now.UTC().Round(time.Microsecond),
	}

	if err := user.HashPassword(); err != nil {
		return store.User{}, err
	}

	const q = `INSERT INTO users(id,email,password,verification_token,verification_expires,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)
	returning id,email,password,verification_token,verification_expires,created_at,updated_at`

	if err := s.db.QueryRowContext(ctx, q, user.ID, user.Email, user.Password, user.VerificationToken, user.VerificationExpires, user.CreatedAt, user.UpdatedAt).Scan(&u.ID, &u.Email, &u.Password, &u.VerificationToken, &u.VerificationExpires, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return store.User{}, errors.Wrap(err, "inserting new user")
	}

	return user, nil
}

// GetUserByID finds a user by id
func (s userStore) GetUserByID(ctx context.Context, id uuid.UUID) (store.User, error) {
	user := store.User{}

	const q = `SELECT id,email,password,verification_token,verified,verification_expires,created_at,updated_at FROM users WHERE id=$1`

	if err := s.db.GetContext(ctx, &user, q, id); err != nil {
		if err == sql.ErrNoRows {
			return store.User{}, errors.New("not found")
		}
		return user, errors.Wrap(err, "retrieving user by id")
	}

	return user, nil
}

// GetUserByEmail retrieves an existing user with the given email.
func (s userStore) GetUserByEmail(ctx context.Context, email string) (store.User, error) {
	user := store.User{}

	const q = `SELECT id,email,password,verification_token,verified,verification_expires,created_at,updated_at FROM users WHERE email=$1`

	if err := s.db.GetContext(ctx, &user, q, email); err != nil {
		return user, errors.Wrap(err, "retrieving user by email")
	}

	return user, nil
}

// GetUserByVerificationToken retrieves user whose verification token matches the given token string.
func (s userStore) GetUserByVerificationToken(ctx context.Context, token uuid.UUID) (store.User, error) {
	user := store.User{}

	const q = `SELECT id,email,password,verification_token,verified,verification_expires,created_at,updated_at FROM users WHERE verification_token=$1`

	if err := s.db.GetContext(ctx, &user, q, token); err != nil {
		return user, errors.Wrap(err, "retrievng user by verification token")
	}

	return user, nil
}

// SetUserAPIKey sets the api key for the given user id.
func (s userStore) SetUserAPIKey(ctx context.Context, id, key uuid.UUID) error {
	user := store.User{
		ID:     id,
		APIKey: &key,
	}

	const q = `UPDATE users SET api_key=$1 WHERE id=$2`

	if _, err := s.db.ExecContext(ctx, q, user.APIKey, user.ID); err != nil {
		return errors.Wrap(err, "updating the api key")
	}

	return nil
}

// SetUserVerified updates the verified value for the user with the given user id.
func (s userStore) SetUserVerified(ctx context.Context, id uuid.UUID) error {
	verified := true

	const q = `UPDATE users SET verified=$1 WHERE id=$2`

	if _, err := s.db.ExecContext(ctx, q, verified, id); err != nil {
		return errors.Wrap(err, "updating verified field")
	}

	return nil
}
