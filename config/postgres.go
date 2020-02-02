package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/go-pg/pg/v9"
)

// PostgresConfig persists the config for our PostgreSQL database connection
type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Database string `env:"POSTGRES_DB"`
}

// GetConnection returns our pg database connection
// usage:
// db := config.GetConnection()
// defer db.Close()
func GetConnection() *pg.DB {
	c := GetPostgresConfig()
	db := pg.Connect(&pg.Options{
		Addr:     c.Host + ":" + c.Port,
		User:     c.User,
		Password: c.Password,
		Database: c.Database,
		PoolSize: 150,
	})
	return db
}

// GetPostgresConfig returns a PostgresConfig pointer with the correct Postgres Config values
func GetPostgresConfig() *PostgresConfig {
	c := PostgresConfig{}
	if err := env.Parse(&c); err != nil {
		fmt.Printf("%+v\n", err)
	}
	return &c
}
