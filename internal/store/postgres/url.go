package postgres

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nairobi-gophers/fupisha/internal/encoding"
	"github.com/nairobi-gophers/fupisha/internal/store"
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
		ID:           encoding.GenUniqueID().String(),
		Owner:        userID,
		OriginalURL:  originalURL,
		ShortenedURL: shortenedURL,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	const q = `INSERT INTO urls (id,owner,original_url,short_url,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6)`

	_, err := u.db.ExecContext(ctx, q, url.ID, url.Owner, url.OriginalURL, url.ShortenedURL, url.CreatedAt, url.UpdatedAt)

	if err != nil {
		return store.URL{}, errors.Wrap(err, "inserting new url")
	}

	return url, nil
}

func (u *urlStore) Get(ctx context.Context, id uuid.UUID) (store.URL, error) {
	var url store.URL

	const q = `SELECT * FROM urls WHERE id=$1`

	_, err := u.db.ExecContext(ctx, q, id, &u)
	if err != nil {
		return store.URL{}, errors.Wrap(err, "retrieving url by id")
	}

	return url, nil
}
