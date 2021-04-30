package url

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"

	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/nairobi-gophers/fupisha/api/v1/auth"
	"github.com/nairobi-gophers/fupisha/encoding"
	"github.com/nairobi-gophers/fupisha/logging"
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

	id, ok := auth.FromContext(r.Context())
	if !ok {
		log(r).Error(errors.New("could not extract userID from context"))
		render.Render(w, r, ErrInternalServerError)
		return
	}

	userID, err := uuid.FromString(id)
	if err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	//Lets validate that userID actually belongs to a real user.
	_, err = rs.Store.Users().Get(r.Context(), userID)
	if err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	link, err := Shorten(body.URL, rs.Config.BaseURL, rs.Config.ParamLength)

	if err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	//Insert the shortened url in the database
	_, err = rs.Store.Urls().New(r.Context(), userID, body.URL, link)
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

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, &resBody)
}

//Shorten shortens a long url string
func Shorten(originalURL, baseURL string, len int) (string, error) {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	param, err := encoding.GenUniqueParam("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890", len)
	if err != nil {
		return "", err
	}

	return baseURL + param, nil
}

func log(r *http.Request) logrus.FieldLogger {
	return logging.GetLogEntry(r)
}
