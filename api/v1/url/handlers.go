package url

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/sirupsen/logrus"

	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/nairobi-gophers/fupisha/api/v1/auth"
	"github.com/nairobi-gophers/fupisha/internal/encoding"
	"github.com/nairobi-gophers/fupisha/internal/logging"
)

type shortenURLRequest struct {
	URL string `json:"url"`
}

func (body *shortenURLRequest) Bind(r *http.Request) error {
	body.URL = strings.TrimSpace(body.URL)

	return validation.ValidateStruct(body, validation.Field(&body.URL, validation.Required, is.URL))
}

//HandleShortenURL shortens the url and returns the shrotened url in the response body
func (rs Resource) HandleShortenURL(w http.ResponseWriter, r *http.Request) {
	body := shortenURLRequest{}

	if err := render.Bind(r, &body); err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	userID, ok := auth.FromContext(r.Context())
	if !ok {
		log(r).Error(errors.New("could not extract userID from context"))
		render.Render(w, r, ErrInternalServerError)
		return
	}

	//Lets validate that userID actually belongs to a real user.
	_, err := rs.Store.Users().Get(userID)
	if err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	link := shortenURL(body.URL, rs.Config.BaseURL, rs.Config.ParamLength)

	//Insert the shortened url in the database
	_, err = rs.Store.Urls().New(userID, body.URL, link)
	if err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	resBody := struct {
		Link string `json:"link"`
	}{
		Link: link,
	}

	render.Respond(w, r, &resBody)
}

func shortenURL(originalURL, baseURL string, len int) string {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	return baseURL + encoding.GenUniqueParam("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890", len)
}

func log(r *http.Request) logrus.FieldLogger {
	return logging.GetLogEntry(r)
}
