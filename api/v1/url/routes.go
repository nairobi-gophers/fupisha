package url

import (
	"github.com/go-chi/chi"
	"github.com/nairobi-gophers/fupisha/api/v1/auth"
)

//Router provides necessary routes for shortening and resolving fupisha urls.
func (rs *Resource) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(auth.Verifier(rs.Config))
		r.Post("/shorten", rs.HandleShortenURL)
	})

	return r
}
