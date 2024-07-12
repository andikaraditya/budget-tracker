package db

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	// MaxConn, IddleConn can be set here.
	// Refer to https://pkg.go.dev/github.com/jackc/pgx/v4@v4.13.0/pgxpool#ParseConfig
	DSN string `envconfig:"dsn" required:"true"`
}

var cfg config

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	if err := envconfig.Process("", &cfg); err != nil {
		fmt.Printf("Error: %v", err)
	}
	fmt.Printf("initialize postgres")
	fmt.Printf("config: %v", cfg)
}
