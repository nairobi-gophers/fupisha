package auth

import "github.com/go-chi/chi"

//Router provides necessary routes for password restricted authentication flow.
func (rs *Resource) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/signup", rs.handleSignup)
	return r
}
