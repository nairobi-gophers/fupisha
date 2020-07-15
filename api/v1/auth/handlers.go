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
)

type signupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (body *signupRequest) Bind(r *http.Request) error {
	body.Email = strings.TrimSpace(body.Email)
	body.Password = strings.TrimSpace(body.Password)
	body.Name = strings.TrimSpace(body.Name)

	return validation.ValidateStruct(body, validation.Field(&body.Email, validation.Required, is.Email), validation.Field(&body.Name, validation.Required, validation.Length(3, 32), is.ASCII),
		validation.Field(&body.Password, validation.Required, validation.Length(8, 32), is.Alphanumeric))
}

//HandleSignup signup handler func for handling requests for new accounts.
func (rs Resource) HandleSignup(w http.ResponseWriter, r *http.Request) {

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

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, http.NoBody)
}

func log(r *http.Request) logrus.FieldLogger {
	return logging.GetLogEntry(r)
}
