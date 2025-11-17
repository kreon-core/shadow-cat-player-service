package config

import (
	"time"

	tul "github.com/kreon-core/shadow-cat-common"
)

type Config struct {
	HTTP    HTTP    `mapstructure:"http"`
	DB      DB      `mapstructure:"db"      validate:"required"`
	Secrets Secrets `mapstructure:"secrets" validate:"required"`
}

type HTTP struct {
	Host              *string        `mapstructure:"host"`
	Port              *int           `mapstructure:"port"`
	ReadTimeout       *time.Duration `mapstructure:"read-timeout"`
	ReadHeaderTimeout *time.Duration `mapstructure:"read-header-timeout"`
	WriteTimeout      *time.Duration `mapstructure:"write-timeout"`
	IdleTimeout       *time.Duration `mapstructure:"idle-timeout"`
}

type DB struct {
	Player tul.PostgresConfig `mapstructure:"player" validate:"required"`
}

type Secrets struct {
	JWTSecretKey string `mapstructure:"jwt-secret-key" validate:"required"`
}
