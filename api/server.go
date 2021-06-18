package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nairobi-gophers/fupisha/config"
	"github.com/nairobi-gophers/fupisha/logging"
	"github.com/nairobi-gophers/fupisha/provider"
)

//Server defines our server dependencies
type Server struct {
	*http.Server
}

//NewServer creates and configures an fupisha API Server serving all application routes.
func NewServer() (*Server, error) {

	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	logger := logging.NewLogger(cfg)

	store, err := cfg.GetStore()

	if err != nil {
		return nil, err
	}

	mailer, err := provider.NewMailerWithSMTP(cfg, "./templates")
	if err != nil {
		return nil, err
	}

	apiCfg := &ApiConfig{
		Logger:     logger,
		Store:      store,
		Cfg:        cfg,
		Mailer:     mailer,
		EnableCORS: false,
	}

	api, err := New(apiCfg)
	if err != nil {
		return nil, err
	}

	srv := http.Server{
		ReadTimeout:  5 * time.Second,   //time from when the connection is accepted to when the request body is fully read
		WriteTimeout: 10 * time.Second,  //time from the end of the request header read to the end of the response write (a.k.a. the lifetime of the ServeHTTP)
		IdleTimeout:  120 * time.Second, // the amount of time a Keep-Alive connection will be kept idle before being reused
		Addr:         ":" + cfg.Port,
		Handler:      api,
	}
	return &Server{&srv}, nil
}

//Start runs ListenAndServe on the http.Server with graceful shutdown.
func (srv *Server) Start() {

	log.Println("Starting Fupisha API Server...")
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()
	log.Printf("Listening on %s\n", srv.Addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit
	log.Println("Shutting down fupisha API server... Reason:", sig)

	//teardown logic here

	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}
	log.Println("Fupisha API server gracefully stopped.")
}
