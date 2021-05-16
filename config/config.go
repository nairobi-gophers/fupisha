package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/nairobi-gophers/fupisha/encoding"
	"github.com/nairobi-gophers/fupisha/store"
	"github.com/nairobi-gophers/fupisha/store/postgres"
)

//Config is a fupisha configuration struct
type Config struct {
	//BaseURL fupisha's fully qualified domain name.
	BaseURL string `envconfig:"FUPISHA_BASE_URL"`
	//Title name of the application e.g. fupisha.
	Title string `envconfig:"FUPISHA_TITLE"`
	//TextLogging write api requests to file.
	TextLogging bool `envconfig:"FUPISHA_TEXT_LOGGING"`
	//LogLevel category of the log.
	LogLevel string `envconfig:"FUPISHA_LOG_LEVEL"`
	//ParamLength length of the shorten url param (https://base_url/{param}) e.g https://fupisha.io/kKIoqRF
	ParamLength int `envconfig:"FUPISHA_PARAM_LENGTH"`
	//Port is the port on which the api server will bind to once started e.g 3333
	Port string `envconfig:"FUPISHA_HTTP_PORT"`
	//JWT json web token payload
	JWT struct {
		//Secret secret jwt signing key.
		Secret string `envconfig:"FUPISHA_JWT_SECRET"`
		//ExpireDelta duration after which the token is rendered invalid.
		ExpireDelta int `envconfig:"FUPISHA_JWT_EXPIRE_DELTA"`
	}
	//SMTP third party email provider smtp configuration fields.
	SMTP struct {
		//Port smtp port
		Port string `envconfig:"FUPISHA_SMTP_PORT"`
		//Host smtp host e.g. smtp.gmail.com
		Host string `envconfig:"FUPISHA_SMTP_HOST"`
		//Username smtp username.
		Username string `envconfig:"FUPISHA_SMTP_USERNAME"`
		//Password smtp password.
		Password string `envconfig:"FUPISHA_SMTP_PASSWORD"`
		//FromName email sender's name.
		FromName string `envconfig:"FUPISHA_SMTP_FROM_NAME"`
		//FromAddress email sender's email address
		FromAddress string `envconfig:"FUPISHA_SMTP_FROM_ADDRESS"`
	}
	//Store fupisha storage configuration object.
	Store struct {
		//Type the type of database. e.g. mongo
		Type string `envconfig:"FUPISHA_STORE_TYPE"`
		//PostgreSQL postgresql database connection parameters.
		PostgreSQL struct {
			//Address postgresql host and port. e.g. localhost:5432
			Address string `envconfig:"FUPISHA_STORE_POSTGRESQL_ADDRESS"`
			//Username postgresql user.
			Username string `envconfig:"FUPISHA_STORE_POSTGRESQL_USERNAME"`
			//Password postgresql password associated with the user.
			Password string `envconfig:"FUPISHA_STORE_POSTGRESQL_PASSWORD"`
			//Database postgresql database name.
			Database string `envconfig:"FUPISHA_STORE_POSTGRESQL_DATABASE"`
			//SSLMode if enabled postgresql will encrypt the communication to and from.
			SSLMode string `envconfig:"FUPISHA_STORE_POSTGRESQL_SSLMODE"`
			//SSLRootCert requires if SSLMode is enabled.
			SSLRootCert string `envconfig:"FUPISHA_STORE_POSTGRESQL_SSLROOTCERT"`
		}
		//Mongo mongo database connection parameters.
		Mongo struct {
			//Address mongo host and port. e.g localhost:27017
			Address string `envconfig:"FUPISHA_STORE_MONGO_ADDRESS"`
			//Username mongo user.
			Username string `envconfig:"FUPISHA_STORE_MONGO_USERNAME"`
			//Password mongo user's password.
			Password string `envconfig:"FUPISHA_STORE_MONGO_PASSWORD"`
			//Database mongo database name.
			Database string `envconfig:"FUPISHA_STORE_MONGO_DATABASE"`
		}
		//MySQL mysql database connection parameters.
		MySQL struct {
			//Address mysql host and port. e.g. localhost:3306
			Address string `envconfig:"FUPISHA_STORE_MYSQL_ADDRESS"`
			//Username mysql user.
			Username string `envconfig:"FUPISHA_STORE_MYSQL_USERNAME"`
			//Password mysql user's passsword.
			Password string `envconfig:"FUPISHA_STORE_MYSQL_PASSWORD"`
			//Database mysql database name.
			Database string `envconfig:"FUPISHA_STORE_MYSQL_DATABASE"`
		}
	}
}

//GetStore returns a connection to the relevant database as specified on the config
func (cfg *Config) GetStore() (store.Store, error) {
	switch cfg.Store.Type {
	case "postgresql":

		dbCfg := &postgres.Config{
			Host:     cfg.Store.PostgreSQL.Address,
			User:     cfg.Store.PostgreSQL.Username,
			Password: cfg.Store.PostgreSQL.Password,
			Name:     cfg.Store.PostgreSQL.Database,
		}

		return postgres.Connect(dbCfg)
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
