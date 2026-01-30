package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port     string `envconfig:"SERVER_PORT" default:"8080"`
	GRPCPort string `envconfig:"GRPC_PORT" default:"50051"`
}

type DatabaseConfig struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     string `envconfig:"DB_PORT" default:"5432"`
	User     string `envconfig:"DB_USER" default:"postgres"`
	Password string `envconfig:"DB_PASSWORD" default:"postgres"`
	Name     string `envconfig:"DB_NAME" default:"employee_db"`
	SSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (d DatabaseConfig) DSN() string {
	return "host=" + d.Host +
		" port=" + d.Port +
		" user=" + d.User +
		" password=" + d.Password +
		" dbname=" + d.Name +
		" sslmode=" + d.SSLMode
}

type JWTConfig struct {
	Secret     string        `envconfig:"JWT_SECRET" default:"your-secret-key-change-in-production"`
	Expiration time.Duration `envconfig:"JWT_EXPIRATION" default:"24h"`
	Issuer     string        `envconfig:"JWT_ISSUER" default:"employee-api"`
}
