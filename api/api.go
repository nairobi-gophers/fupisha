package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/nairobi-gophers/fupisha/api/v1/auth"
	"github.com/nairobi-gophers/fupisha/internal/config"
	"github.com/nairobi-gophers/fupisha/internal/logging"
	"github.com/nairobi-gophers/fupisha/internal/store"
)

//New configures application resources and routers.
func New(enableCORS bool, cfg *config.Config, store store.Store) (*chi.Mux, error) {
	logger := logging.NewLogger(cfg)

	authResource := auth.NewResource(store, cfg)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(logging.NewStructuredLogger(logger))
	r.Use(auth.CheckAPI)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	//Use CORS middleware if client is not served by this api, e.g. from other domain
	//or CDN
	if enableCORS {
		r.Use(corsConfig().Handler)
	}

	r.Mount("/auth", authResource.Router())

	r.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "User-agent: *\nDisallow: /")
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	fmt.Println("[+] API ROUTES ")

	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("walkFunc err:%s\n", err.Error())
	}

	fmt.Println("[-] API ROUTES")

	return r, nil
}

func corsConfig() *cors.Cors {
	return cors.New(cors.Options{

		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Accept-Encoding", "Content-Type", "Content-Length", "X-CSRF-Token"},
		ExposedHeaders:     []string{"Link"},
		AllowCredentials:   true,
		MaxAge:             86400,
		OptionsPassthrough: false,
	})
}
