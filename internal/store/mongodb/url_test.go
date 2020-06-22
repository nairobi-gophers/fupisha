package mongodb

import (
	"reflect"
	"testing"
	"time"

	"github.com/nairobi-gophers/fupisha/internal/store/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUrl(t *testing.T) {
	s, teardown := testConn(t)
	defer teardown("urls")

	id, err := s.Urls().New("5d0575344d9f7ff15e989174", "https:///example.com/shorten-me/2016/caching/6/caching-with-golang/", "https://fupi.sha/sDY6KJ")

	if err != nil {
		t.Fatalf("failed to create url %s: ", err)
	}

	if _, ok := id.(primitive.ObjectID); !ok {
		t.Fatalf("failed to assert the created url insert id")
	}

	uid := id.(primitive.ObjectID).Hex()

	url, err := s.Urls().Get(uid)

	if err != nil {
		t.Fatalf("failed to get url by id: %s", err)
	}

	sinceCreated := time.Since(url.CreatedAt)
	if sinceCreated > 2*time.Second || sinceCreated < 0 {
		t.Fatalf("bad url.CreatedAt: %v", url.CreatedAt)
	}

	userID, err := primitive.ObjectIDFromHex("5d0575344d9f7ff15e989174")
	if err != nil {
		t.Fatalf("failed to extract object id from string: %s", err)
	}

	want := model.URL{
		ID:           id.(primitive.ObjectID),
		User:         userID,
		OriginalURL:  "https:///example.com/shorten-me/2016/caching/6/caching-with-golang/",
		ShortenedURL: "https://fupi.sha/sDY6KJ",
		CreatedAt:    url.CreatedAt,
		UpdatedAt:    time.Time{},
	}

	if !reflect.DeepEqual(url, want) {
		t.Fatalf("got url %+v want %+v", url, want)
	}

}
