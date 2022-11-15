package url

import (
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
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

// HandleShortenURL shortens the url and returns the shrotened url in the response body
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
		log(r).WithField("id", id).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	//Lets validate that userID actually belongs to a real user.
	_, err = rs.Store.GetUserByID(r.Context(), userID)
	if err != nil {
		log(r).WithField("userID", userID).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	baseURL := rs.Config.BaseURL + ":" + rs.Config.Port

	param, err := Shorten(body.URL, rs.Config.ParamLength)

	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	link := baseURL + param

	if err != nil {
		log(r).WithField("url", body.URL).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	type resBody struct {
		Link string `json:"link"`
	}

	//Insert the shortened url in the database
	_, err = rs.Store.NewURL(r.Context(), userID, body.URL, param)
	if err != nil {
		if pqErr, ok := errors.Cause(err).(*pq.Error); ok {
			//if its a unique key violation, that means we had already shortened the url before.
			if pqErr.Code == pq.ErrorCode("23505") {
				//Let's retrieve the shortened url param.
				url, err := rs.Store.GetURLByLongStr(r.Context(), body.URL)
				if err != nil {
					log(r).WithField("url", body.URL).Error(err)
					render.Render(w, r, ErrInternalServerError)
					return
				}
				//concatenate the short url param with our baseurl e.g
				//http://localhost:8888/ + okzbUwy = http://localhost:8888/okzbUwy
				link = baseURL + url.ShortenedURLParam

				resp := resBody{
					Link: link,
				}

				render.Status(r, http.StatusCreated)
				render.Respond(w, r, &resp)
				return
			}
		}
		log(r).WithField("url", body.URL).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	resp := resBody{
		Link: link,
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, &resp)
}

// Shorten shortens a long url string
func Shorten(originalURL string, len int) (string, error) {

	param, err := encoding.GenUniqueParam("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890", len)
	if err != nil {
		return "", err
	}

	return param, nil
}

func log(r *http.Request) logrus.FieldLogger {
	return logging.GetLogEntry(r)
}
