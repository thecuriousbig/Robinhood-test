package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	App      app
	Endpoint endpoint
	Mongo    mongo
	JWT      jwt
}

type app struct {
	EnableSwagger bool `envconfig:"ENABLE_SWAGGER" default:"true"`
}

type endpoint struct {
	Port string `envconfig:"PORT" default:"8080"`
}

type mongo struct {
	URI      string `envconfig:"MONGO_URI" default:"mongodb://localhost:27017"`
	Database string `envconfig:"DB_NAME" default:"robinhood"`
}

type jwt struct {
	Secret          string `envconfig:"JWT_SECRET"`
	AUD             string `envconfig:"JWT_AUD"`
	ISS             string `envconfig:"JWT_ISS"`
	ExpiresHours    uint   `envconfig:"JWT_EXPIRES_HOURS" default:"730"`
	AutoLogoffHours uint   `envconfig:"JWT_AUTO_LOGOFF_HOURS" default:"730"`
}

var cfg config

func New() {
	_ = godotenv.Load()
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("read env error : %s", err.Error())
	}
}

func Get() config {
	return cfg
}
