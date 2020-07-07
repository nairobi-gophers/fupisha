package config

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/nairobi-gophers/fupisha/internal/encoding"
	"github.com/nairobi-gophers/fupisha/internal/store"
	"github.com/nairobi-gophers/fupisha/internal/store/mongodb"
)

//Config is a fupisha configuration struct
type Config struct {
	BaseURL     string `envconfig:"FUPISHA_BASE_URL"`
	Title       string `envconfig:"FUPISHA_TITLE"`
	TextLogging bool   `envconfig:"FUPISHA_TEXT_LOGGING"`
	LogLevel    string `envconfig:"FUPISHA_LOG_LEVEL"`
	JWT         struct {
		Secret      string `envconfig:"FUPISHA_JWT_SECRET"`
		ExpireDelta string `envconfig:"FUPISHA_JWT_EXPIRE_DELTA"`
	}

	SMTP struct {
		Port        int    `envconfig:"FUPISHA_SMTP_PORT"`
		Host        string `envconfig:"FUPISHA_SMTP_HOST"`
		Username    string `envconfig:"FUPISHA_SMTP_USERNAME"`
		Password    string `envconfig:"FUPISHA_SMTP_PASSWORD"`
		FromName    string `envconfig:"FUPISHA_SMTP_FROM_NAME"`
		FromAddress string `envconfig:"FUPISHA_SMTP_FROM_ADDRESS"`
	}

	Store struct {
		Type string `envconfig:"FUPISHA_STORE_TYPE"`

		PostgreSQL struct {
			Address     string `envconfig:"FUPISHA_STORE_POSTGRESQL_ADDRESS"`
			Username    string `envconfig:"FUPISHA_STORE_POSTGRESQL_USERNAME"`
			Password    string `envconfig:"FUPISHA_STORE_POSTGRESQL_PASSWORD"`
			Database    string `envconfig:"FUPISHA_STORE_POSTGRESQL_DATABASE"`
			SSLMode     string `envconfig:"FUPISHA_STORE_POSTGRESQL_SSLMODE"`
			SSLRootCert string `envconfig:"FUPISHA_STORE_POSTGRESQL_SSLROOTCERT"`
		}

		Mongo struct {
			Address  string `envconfig:"FUPISHA_STORE_MONGO_ADDRESS"`
			Username string `envconfig:"FUPISHA_STORE_MONGO_USERNAME"`
			Password string `envconfig:"FUPISHA_STORE_MONGO_PASSWORD"`
			Database string `envconfig:"FUPISHA_STORE_MONGO_DATABASE"`
		}

		MySQL struct {
			Address  string `envconfig:"FUPISHA_STORE_MYSQL_ADDRESS"`
			Username string `envconfig:"FUPISHA_STORE_MYSQL_USERNAME"`
			Password string `envconfig:"FUPISHA_STORE_MYSQL_PASSWORD"`
			Database string `envconfig:"FUPISHA_STORE_MYSQL_DATABASE"`
		}
	}
}

//GetStore returns a connection to the relevant database as specified on the config
func (cfg *Config) GetStore() (store.Store, error) {
	switch cfg.Store.Type {
	case "mongo":
		var (
			address  = cfg.Store.Mongo.Address
			username = cfg.Store.Mongo.Username
			password = cfg.Store.Mongo.Password
			database = cfg.Store.Mongo.Database
		)

		if (address == "") || (username == "") || (database == "") || (password == "") {
			err := errors.New("missing mongodb configuration variable")
			log.Fatalf("config: %s", err.Error())
		}
		// address := fmt.Sprintf("%s:%s", host, port)
		return mongodb.Connect(
			address,
			username,
			password,
			database,
		)
	}
	return nil, fmt.Errorf("config: unknown store type: %s", cfg.Store.Type)
}

//New returns an initialized config object ready for use
func New() (*Config, error) {
	cfg := Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to process environment variables: %v", err)
	}
	cfg.BaseURL = strings.TrimSuffix(cfg.BaseURL, "/")

	return &cfg, nil
}

//GenKey generates a  32 byte crypto-random unique key
func GenKey() {
	log.Printf("%s\n", encoding.GenHexKey(32))
}
