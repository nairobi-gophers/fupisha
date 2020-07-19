package auth

import (
	"encoding/hex"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/nairobi-gophers/fupisha/internal/logging"
	"github.com/nairobi-gophers/fupisha/internal/provider"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type signupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
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

func (body *loginRequest) Bind(r *http.Request) error {
	body.Email = strings.TrimSpace(body.Email)
	body.Password = strings.TrimSpace(body.Password)

	return validation.ValidateStruct(body, validation.Field(&body.Email, validation.Required, is.Email), validation.Field(&body.Password, validation.Required, validation.Length(8, 32), is.Alphanumeric))
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

//HandleLogin login handler for handling login requests
func (rs Resource) HandleLogin(w http.ResponseWriter, r *http.Request) {

	body := loginRequest{}

	if err := render.Bind(r, &body); err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	usr, err := rs.Store.Users().GetByEmail(body.Email)
	if err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrUnauthorized(ErrInvalidEmailOrPassword))
		return
	}

	if _, err := usr.Compare(usr.Password, body.Password); err != nil {
		log(r).WithField("password", body.Password).Error(err)
		render.Render(w, r, ErrUnauthorized(ErrInvalidEmailOrPassword))
		return
	}

	secret := hex.EncodeToString([]byte(rs.Config.JWT.Secret))
	jwtService, err := provider.NewJWTService(secret)
	if err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	token, err := jwtService.Encode(usr.ID.Hex())
	if err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	resBody := struct {
		Email  string `json:"email"`
		Name   string `json:"name"`
		UserID string `json:"id"`
		Token  string `json:"token"`
	}{
		Email:  usr.Email,
		Name:   usr.Name,
		UserID: usr.ID.Hex(),
		Token:  token,
	}

	render.Respond(w, r, &resBody)
}

func log(r *http.Request) logrus.FieldLogger {
	return logging.GetLogEntry(r)
}
