package auth

import (
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/nairobi-gophers/fupisha/internal/auth"
	"github.com/nairobi-gophers/fupisha/internal/logging"
)

//Verifier http middleware will verify a jwt string from a http request.
func (rs *Resource) Verifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] != nil {
			authHeader := r.Header.Get("Authorization")
			if len(authHeader) < 7 || strings.ToUpper(authHeader[:6]) != "BEARER" {
				logging.FromContext(r.Context()).Error(ErrLoginToken)
				render.Render(w, r, ErrUnauthorized(ErrLoginToken))
				return
			}
			token := authHeader[7:]

			secret := hex.EncodeToString([]byte(rs.Config.JWT.Secret))
			service, err := auth.NewService(secret)
			if err != nil {
				logging.FromContext(r.Context()).Error(err)
				render.Render(w, r, ErrInternalServerError)
				return
			}
			_, _, err = service.Decode(token)
			if err != nil {
				logging.FromContext(r.Context()).Error(err)
				render.Render(w, r, ErrUnauthorized(ErrLoginToken))
				return
			}
		} else {
			logging.FromContext(r.Context()).Error(ErrMissingToken)
			render.Render(w, r, ErrUnauthorized(ErrMissingToken))
			return
		}
	})
}
