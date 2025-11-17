package config

import (
	"time"

	"github.com/kreon-core/shadow-cat-common/postgres"
)

type Config struct {
	HTTP HTTP `mapstructure:"http"`
	DB   DB   `mapstructure:"db"   validate:"required"`
}

type HTTP struct {
	Host         *string        `mapstructure:"host"`
	Port         *int           `mapstructure:"port"`
	ReadTimeout  *time.Duration `mapstructure:"read_timeout"`
	WriteTimeout *time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  *time.Duration `mapstructure:"idle_timeout"`
}

type DB struct {
	Player postgres.Config `mapstructure:"player" validate:"required"`
}
