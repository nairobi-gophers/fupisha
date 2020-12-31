package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/nairobi-gophers/fupisha/internal/encoding"
	"github.com/nairobi-gophers/fupisha/internal/store/model"
	"github.com/pkg/errors"
)

type userStore struct {
	db *sqlx.DB
}

//New creates a new user record.
func (s userStore) New(ctx context.Context, name, email, password string) (model.User, error) {
	tkn := encoding.GenUniqueID()

	var now time.Time = time.Now()

	user := model.User{
		ID:                  encoding.GenUniqueID().String(),
		Name:                name,
		Email:               email,
		Password:            password,
		VerificationToken:   tkn,
		VerificationExpires: now.Add(time.Minute * 60), //expires 60 mins later
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if err := user.HashPassword(); err != nil {
		return model.User{}, err
	}

	const q = `INSERT INTO users(id,name,email,password,verification_token,verification_expires,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	if _, err := s.db.ExecContext(ctx, q, user.ID, user.Name, user.Email, user.Password, user.VerificationToken, user.VerificationExpires, user.CreatedAt, user.UpdatedAt); err != nil {
		return model.User{}, errors.Wrap(err, "inserting new user")
	}

	return user, nil
}

//Get finds a user by id
func (s userStore) Get(ctx context.Context, id string) (model.User, error) {
	user := model.User{}

	const q = `SELECT id,name,email,api_key,created_at,verified FROM users WHERE id=$1`

	if err := s.db.GetContext(ctx, &user, q, id); err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, errors.New("not found")
		}
		return user, errors.Wrap(err, "retrieving user by id")
	}

	return user, nil
}

//GetByEmail retrieves an existing user with the given email.
func (s userStore) GetByEmail(ctx context.Context, email string) (model.User, error) {
	user := model.User{}

	const q = `SELECT id,name,email,api_key,created_at,verified FROM users WHERE email=$1`

	if err := s.db.GetContext(ctx, &user, q, email); err != nil {
		return user, errors.Wrap(err, "retrieving user by email")
	}

	return user, nil
}

//SetAPIKey sets the api key for the given user id.
func (s userStore) SetAPIKey(ctx context.Context, id string, key uuid.UUID) error {
	user := model.User{
		ID:     id,
		APIKey: key,
	}

	const q = `UPDATE users SET api_key=$1 WHERE id=$2`

	if _, err := s.db.ExecContext(ctx, q, user.APIKey, user.ID); err != nil {
		return errors.Wrap(err, "updating the api key")
	}

	return nil
}
