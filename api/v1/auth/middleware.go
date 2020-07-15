package auth

import (
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/nairobi-gophers/fupisha/internal/auth"
)

const (
	//version of api provided by the server
	apiVersion = "v1"
)

//Verifier http middleware will verify a jwt string from a http request.
func (rs *Resource) Verifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] != nil {
			authHeader := r.Header.Get("Authorization")
			if len(authHeader) < 7 || strings.ToUpper(authHeader[:6]) != "BEARER" {
				log(r).Error(ErrLoginToken)
				render.Render(w, r, ErrUnauthorized(ErrLoginToken))
				return
			}
			token := authHeader[7:]

			secret := hex.EncodeToString([]byte(rs.Config.JWT.Secret))
			service, err := auth.NewService(secret)
			if err != nil {
				log(r).WithField("secret", secret).Error(err)
				render.Render(w, r, ErrInternalServerError)
				return
			}
			_, _, err = service.Decode(token)
			if err != nil {
				log(r).WithField("token", token).Error(err)
				render.Render(w, r, ErrUnauthorized(ErrLoginToken))
				return
			}
		} else {
			log(r).Error(ErrMissingToken)
			render.Render(w, r, ErrUnauthorized(ErrMissingToken))
			return
		}
		next.ServeHTTP(w, r)
	})
}

//CheckAPI http middleware will verify the api version from a http request.
func CheckAPI(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Api"] != nil {
			api := r.Header.Get("Api")
			if len(api) > 0 {
				if apiVersion != api {
					log(r).Error(ErrAPIUnsupported)
					render.Render(w, r, ErrUnsupportedAPIVersion(ErrAPIUnsupported))
					return
				}
			} else {
				log(r).Error(ErrAPIUnsupported)
				render.Render(w, r, ErrUnsupportedAPIVersion(ErrAPIUnsupported))
				return
			}
		} else {
			log(r).Error(ErrMissingAPIVersion)
			render.Render(w, r, ErrUnsupportedAPIVersion(ErrMissingAPIVersion))
			return
		}
		next.ServeHTTP(w, r)
	})
}
