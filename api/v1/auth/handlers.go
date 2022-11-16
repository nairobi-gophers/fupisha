package auth

import (
	"net/http"
	"strings"
	"time"

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

// HandleSignup signup handler func for handling requests for new accounts.
func (rs Resource) HandleSignup(w http.ResponseWriter, r *http.Request) {
	body := signupRequest{}

	if err := render.Bind(r, &body); err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	u, err := rs.Store.NewUser(r.Context(), body.Email, body.Password)

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
		VerificationURL:    rs.Config.BaseURL + ":" + rs.Config.Port + "/auth/verify/?v=" + encoding.Encode(u.VerificationToken),
	}

	errc := make(chan error, 1)
	go func(err chan error) {
		err <- rs.Mailer.SendVerifyNotification(body.Email, verifyEmailContent)
	}(errc)

	if err := <-errc; err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	resBody := struct {
		Status string `json:"status"`
		Data   string `json:"data"`
	}{
		Status: http.StatusText(http.StatusOK),
		Data:   "signup successful, check your email to verify your account",
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, &resBody)
}

// HandleVerify verify the verification code sent with the signup email
func (rs Resource) HandleVerify(w http.ResponseWriter, r *http.Request) {
	verificationCode := r.URL.Query().Get("v")

	//verification token should not be empty
	if len(verificationCode) == 0 {
		log(r).WithField("verificationcode", verificationCode)
		render.Render(w, r, ErrInvalidRequest(ErrInvalidVerificationToken))
		return
	}

	//we decode the verification code
	code, err := encoding.Decode(verificationCode)
	if err != nil {
		log(r).WithField("verificationcode", verificationCode)
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	//then check if it exists in the database
	u, err := rs.Store.GetUserByVerificationToken(r.Context(), code)
	if err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInvalidRequest(ErrInvalidVerificationToken))
		return
	}

	//check if token is expired
	if time.Until(u.VerificationExpires) < 1 {
		log(r).WithField("verificationtoken", verificationCode)
		render.Render(w, r, ErrInvalidRequest(ErrInvalidVerificationToken))
		return
	}

	//mark the user as verified
	if err := rs.Store.SetUserVerified(r.Context(), u.ID); err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	//send a welcome email.
	welcomeEmailContent := provider.WelcomeEmailContent{
		LoginURL: "https://fupisha.io/login",
		SiteName: "Fupisha",
		SiteURL:  "https://fupisha.io",
	}

	errc := make(chan error, 1)
	go func(err chan error) {
		err <- rs.Mailer.SendWelcomeNotification(u.Email, welcomeEmailContent)
	}(errc)

	if err := <-errc; err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	resBody := struct {
		Status string `json:"status"`
		Data   string `json:"data"`
	}{
		Status: http.StatusText(http.StatusOK),
		Data:   "verification successful.",
	}

	//We should redirect to a frontend page once the frontend is up and running. The page should have a text probably saying account verified successfully and that the user should have recieved a welcome email with
	//login instructions.
	render.Status(r, http.StatusOK)
	render.Respond(w, r, &resBody)
}

// HandleLogin login handler for handling login requests
func (rs Resource) HandleLogin(w http.ResponseWriter, r *http.Request) {

	body := loginRequest{}

	if err := render.Bind(r, &body); err != nil {
		log(r).Error(err)
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	usr, err := rs.Store.GetUserByEmail(r.Context(), body.Email)
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
