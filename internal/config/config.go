package config

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

