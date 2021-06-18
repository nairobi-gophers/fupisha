package auth

import (
	"net/http"
	"strings"

	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/lib/pq"
	"github.com/nairobi-gophers/fupisha/encoding"
	"github.com/nairobi-gophers/fupisha/logging"
	"github.com/nairobi-gophers/fupisha/provider"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type signupRequest struct {
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

	return validation.ValidateStruct(body, validation.Field(&body.Email, validation.Required, is.Email), validation.Field(&body.Password, validation.Required, validation.Length(8, 32), is.Alphanumeric))
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

	u, err := rs.Store.Users().New(r.Context(), body.Email, body.Password)

	if err != nil {
		if pqErr, ok := errors.Cause(err).(*pq.Error); ok {
			//If its a unique key violation
			if pqErr.Code == pq.ErrorCode("23505") {
				log(r).WithField("email", body.Email).Error(err)
				render.Render(w, r, ErrDuplicateField(ErrEmailTaken))
				return
			}
		}
		log(r).WithField("email", body.Email).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	verifyEmailContent := provider.VerifyEmailContent{
		SiteURL:            "http://fupisha.io",
		SiteName:           "Fupisha",
		VerificationExpiry: u.VerificationExpires,
		VerificationURL:    rs.Config.BaseURL + rs.Config.Port + "/verify/?v=" + encoding.Encode(u.VerificationToken),
	}

	if err := rs.Mailer.SendVerifyNotification(body.Email, verifyEmailContent); err != nil {
		log(r).Error(err)
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

	usr, err := rs.Store.Users().GetByEmail(r.Context(), body.Email)
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

	jwtService, err := provider.NewJWTService(rs.Config)
	if err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	token, err := jwtService.Encode(usr.ID.String())
	if err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	resBody := struct {
		Email  string `json:"email"`
		UserID string `json:"id"`
		Token  string `json:"token"`
	}{
		Email:  usr.Email,
		UserID: encoding.Encode(usr.ID),
		Token:  token,
	}

	render.Respond(w, r, &resBody)
}

func log(r *http.Request) logrus.FieldLogger {
	return logging.GetLogEntry(r)
}
