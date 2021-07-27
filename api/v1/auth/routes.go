package auth

import "github.com/go-chi/chi"

//Router provides necessary routes for password restricted authentication flow.
func (rs *Resource) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/verify", rs.HandleVerify) //verify verification token using url query i.e /auth/verify?v="assgsggsghahs563782"
	r.Group(func(r chi.Router) {
		r.Use(CheckAPI)
		r.Post("/signup", rs.HandleSignup)
		r.Post("/login", rs.HandleLogin)
	})
	return r
}
