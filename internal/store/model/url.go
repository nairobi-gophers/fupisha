package model

import (
	"time"
)

//URL contains all the related info about the shortened url.
type URL struct {
	ID           string    `bson:"_id,omitempty" db:"id"`
	User         string    `bson:"user,omitempty" db:"user"`
	OriginalURL  string    `bson:"originalURL,omitempty" db:"original_url"`
	ShortenedURL string    `bson:"shortenedURL,omitempty" db:"short_url"`
	VisitCount   int       `bson:"visitCount,omitempty" db:"visit_count"`
	CreatedAt    time.Time `bson:"createdAt,omitempty" db:"created_at"`
	UpdatedAt    time.Time `bson:"updatedAt,omitempty" db:"updated_at"`
}
