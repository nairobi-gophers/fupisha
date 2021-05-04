package postgres

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nairobi-gophers/fupisha/encoding"
	"github.com/nairobi-gophers/fupisha/store"
	"github.com/pkg/errors"
)

type urlStore struct {
	db *sqlx.DB
}

//New created a new url record.
func (u *urlStore) New(ctx context.Context, userID uuid.UUID, originalURL, shortenedURL string) (store.URL, error) {

	//Lets check if its a valid UUID
	// if _, err := uuid.FromString(userID); err != nil {
	// 	return store.URL{}, errors.Wrap(err, "invalid uuid userID")
	// }

	now := time.Now()

	url := store.URL{
		ID:           encoding.GenUniqueID(),
		Owner:        userID,
		OriginalURL:  originalURL,
		ShortenedURL: shortenedURL,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	var ur store.URL

	const q = `INSERT INTO urls (id,owner,original_url,short_url,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6) returning id,owner,original_url,short_url,created_at,updated_at`

	if err := u.db.QueryRowContext(ctx, q, url.ID, url.Owner, url.OriginalURL, url.ShortenedURL, url.CreatedAt, url.UpdatedAt).Scan(&ur.ID, &ur.Owner, &ur.OriginalURL, &ur.ShortenedURL, &ur.CreatedAt, &ur.UpdatedAt); err != nil {
		return store.URL{}, errors.Wrap(err, "inserting new url")
	}

	return ur, nil
}

func (u *urlStore) Get(ctx context.Context, id uuid.UUID) (store.URL, error) {
	var url store.URL

	const q = `SELECT * FROM urls WHERE id=$1`

	if err := u.db.GetContext(ctx, &url, q, id); err != nil {
		return store.URL{}, errors.Wrap(err, "retrieving url by id")
	}

	return url, nil
}

func (u *urlStore) GetByParam(ctx context.Context, param string) (store.URL, error) {
	var url store.URL

	const q = `SELECT * FROM urls WHERE short_url=$1`
	if err := u.db.GetContext(ctx, &url, q, param); err != nil {
		return store.URL{}, errors.Wrap(err, "retrieving url by param")
	}

	return url, nil
}
