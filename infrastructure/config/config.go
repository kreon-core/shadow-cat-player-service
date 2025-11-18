package config

import (
	"time"

	"github.com/kreon-core/shadow-cat-common/dbc"
)

type Config struct {
	HTTP      HTTP      `mapstructure:"http"      validate:"required"`
	Databases Databases `mapstructure:"databases" validate:"required"`
	Externals Externals `mapstructure:"externals" validate:"required"`
}

type HTTP struct {
	Host              string         `mapstructure:"host"                validate:"required"`
	Port              int            `mapstructure:"port"                validate:"required"`
	ReadTimeout       *time.Duration `mapstructure:"read-timeout"`
	ReadHeaderTimeout *time.Duration `mapstructure:"read-header-timeout"`
	WriteTimeout      *time.Duration `mapstructure:"write-timeout"`
	IdleTimeout       *time.Duration `mapstructure:"idle-timeout"`
}

type Databases struct {
	Player dbc.PostgresConfig `mapstructure:"player" validate:"required"`
}
type Externals struct {
	AuthClient Client `mapstructure:"auth-client" validate:"required"`
}
type Client struct {
	BaseURL string            `mapstructure:"base-url" validate:"required,url"`
	Paths   map[string]string `mapstructure:"paths"    validate:"required"`
}
