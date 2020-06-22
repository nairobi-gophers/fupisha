package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//URL contains all the related info about the shortened url.
type URL struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	User         primitive.ObjectID `bson:"user,omitempty"`
	OriginalURL  string             `bson:"originalURL,omitempty"`
	ShortenedURL string             `bson:"shortenedURL,omitempty"`
	VisitCount   int                `bson:"visitCount,omitempty"`
	Target       string             `bson:"target,omitempty"`
	CreatedAt    time.Time          `bson:"createdAt,omitempty"`
	UpdatedAt    time.Time          `bson:"updatedAt,omitempty"`
}
