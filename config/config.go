package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config
type Config struct {
	RESTServer RESTServer
	GRPCServer GRPCServer
	Postgres   Postgres
}

// RESTServer config
type RESTServer struct {
	RestPort string `env:"REST_PORT"`
}

// RESTServer config
type GRPCServer struct {
	GrpcPort string `env:"GRPC_PORT"`
}

// Postgresql config
type Postgres struct {
	PostgresqlHost     string `env:"POSTGRES_HOST"`
	PostgresqlUser     string `env:"POSTGRES_USER"`
	PostgresqlPassword string `env:"POSTGRES_PASSWORD"`
	PostgresqlDbname   string `env:"POSTGRES_DB"`
}

var (
	config *Config
	once   sync.Once
)

// Get the config file
func GetConfig() *Config {
	once.Do(func() {
		log.Println("read application configuration")
		config = &Config{}
		if err := cleanenv.ReadConfig(".env", config); err != nil {
			help, _ := cleanenv.GetDescription(config, nil)
			log.Println(help)
			log.Fatal(err)
		}
	})
	return config
}
