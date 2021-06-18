package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/nairobi-gophers/fupisha/api/v1/auth"
	"github.com/nairobi-gophers/fupisha/api/v1/url"
	"github.com/nairobi-gophers/fupisha/config"
	"github.com/nairobi-gophers/fupisha/logging"
	"github.com/nairobi-gophers/fupisha/provider"
	"github.com/nairobi-gophers/fupisha/store"
	"github.com/sirupsen/logrus"
)

//ApiConfig declares the required api server dependencies.
type ApiConfig struct {
	Logger     *logrus.Logger
	Cfg        *config.Config
	Store      store.Store
	Mailer     *provider.Mailer
	EnableCORS bool
}

//New configures application resources and routers.
func New(apiCfg *ApiConfig) (*chi.Mux, error) {

	authResource := auth.NewResource(apiCfg.Store, apiCfg.Cfg, apiCfg.Mailer)
	urlResource := url.NewResource(apiCfg.Store, apiCfg.Cfg)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(logging.NewStructuredLogger(apiCfg.Logger))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	//Use CORS middleware if client is not served by this api, e.g. from other domain
	//or CDN
	if apiCfg.EnableCORS {
		r.Use(corsConfig().Handler)
	}

	r.Mount("/auth", authResource.Router())
	r.Mount("/url", urlResource.Router())

	//Redirect shortened urls
	r.Get("/{urlParam}", func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "urlParam")
		u, err := apiCfg.Store.Urls().GetByParam(r.Context(), param)
		if err != nil {
			logging.GetLogEntry(r).WithField("param", param).Error(err)
			render.Render(w, r, url.ErrURLNotFound(errors.New("not found")))
			return
		}
		http.Redirect(w, r, u.OriginalURL, http.StatusFound)
	})

	r.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "User-agent: *\nDisallow: /")
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
	// 	route = strings.Replace(route, "/*/", "/", -1)
	// 	fmt.Printf("%s %s\n", method, route)
	// 	return nil
	// }

	// fmt.Println("[+] API ROUTES ")

	// if err := chi.Walk(r, walkFunc); err != nil {
	// 	fmt.Printf("walkFunc err:%s\n", err.Error())
	// }

	// fmt.Println("[-] API ROUTES")

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
