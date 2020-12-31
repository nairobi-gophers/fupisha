package model

import (
	"time"
)

//URL contains all the related info about the shortened url.
type URL struct {
	ID           string    `bson:"_id,omitempty" db:"id"`
	User         string    `bson:"user,omitempty" db:"user"`
	OriginalURL  string    `bson:"originalURL,omitempty" db:"originalURL"`
	ShortenedURL string    `bson:"shortenedURL,omitempty" db:"shortenedURL"`
	VisitCount   int       `bson:"visitCount,omitempty" db:"visitCount"`
	CreatedAt    time.Time `bson:"createdAt,omitempty" db:"createdAt"`
	UpdatedAt    time.Time `bson:"updatedAt,omitempty" db:"updatedAt"`
}
