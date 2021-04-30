package store

import (
	"time"

	"github.com/gofrs/uuid"
)

//URL contains all the related info about the shortened url.
type URL struct {
	ID           string    `db:"id"`
	Owner        uuid.UUID `db:"owner"`
	OriginalURL  string    `db:"original_url"`
	ShortenedURL string    `db:"short_url"`
	VisitCount   *int      `db:"visit_count,omitempty"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
