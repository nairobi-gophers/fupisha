package url

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/nairobi-gophers/fupisha/internal/encoding"
	"github.com/nairobi-gophers/fupisha/internal/logging"
)

func shortenURL(originalURL, baseURL string, len int) string {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	return baseURL + encoding.GenUniqueParam("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890", len)
}

func log(r *http.Request) logrus.FieldLogger {
	return logging.GetLogEntry(r)
}
