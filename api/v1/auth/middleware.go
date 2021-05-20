package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/nairobi-gophers/fupisha/config"
	"github.com/nairobi-gophers/fupisha/provider"
)

//The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key int

const (
	//version of api provided by the server
	apiVersion = "v1"
	//userIDKey for propagating the userID down the request chain.
	userIDKey key = iota //use of named type to stop golint from complaining.
	//https://blog.golang.org/context#TOC_3.2.
)

//Verifier http middleware will verify a jwt string from a http request.
func Verifier(cfg *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			if r.Header["Authorization"] != nil {
				authHeader := r.Header.Get("Authorization")
				if len(authHeader) < 7 || strings.ToUpper(authHeader[:6]) != "BEARER" {
					log(r).Error(ErrLoginToken)
					render.Render(w, r, ErrUnauthorized(ErrLoginToken))
					return
				}
				token := authHeader[7:]

				service, err := provider.NewJWTService(cfg)
				if err != nil {
					log(r).WithField("secret", cfg.JWT.Secret).Error(err)
					render.Render(w, r, ErrInternalServerError)
					return
				}
				uid, _, err := service.Decode(token)
				if err != nil {
					log(r).WithField("token", token).Error(err)
					render.Render(w, r, ErrUnauthorized(ErrLoginToken))
					return
				}

				ctx := context.WithValue(r.Context(), userIDKey, uid)

				next.ServeHTTP(w, r.WithContext(ctx))

			} else {
				log(r).Error(ErrMissingToken)
				render.Render(w, r, ErrUnauthorized(ErrMissingToken))
				return
			}
		}
		return http.HandlerFunc(hfn)
	}
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
