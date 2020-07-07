package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/nairobi-gophers/fupisha/internal/logging"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type signupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (body *signupRequest) Bind(r *http.Request) error {
	body.Email = strings.TrimSpace(body.Email)
	body.Password = strings.TrimSpace(body.Password)

	return validation.ValidateStruct(body, validation.Field(&body.Email, validation.Required, is.Email),
		validation.Field(&body.Password, validation.Required, validation.Length(8, 32), is.Alphanumeric))
}

func (rs Resource) handleSignup(w http.ResponseWriter, r *http.Request) {

	if err := checkAPI(r.Header.Get("Api")); err != nil {
		log(r).WithField("APIVersion", r.Header.Get("Api")).Error(err)
		render.Render(w, r, ErrUnsupportedAPIVersion(err))
		return
	}

	body := signupRequest{}

	if err := render.Bind(r, &body); err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	_, err := rs.Store.Users().New(body.Name, body.Email, body.Password)

	if err != nil {
		var e mongo.WriteException
		if errors.As(err, &e) {
			//If its a unique key violation
			for _, we := range e.WriteErrors {
				if we.Code == 11000 {
					log(r).WithField("email", body.Email).Error(err)
					render.Render(w, r, ErrDuplicateField(ErrEmailTaken))
					return
				}
			}
		}
		log(r).WithField("email", body.Email).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.Respond(w, r, http.NoBody)
}

func checkAPI(api string) error {
	//if  API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

func log(r *http.Request) logrus.FieldLogger {
	return logging.GetLogEntry(r)
}
