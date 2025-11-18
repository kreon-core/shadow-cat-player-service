package infrastructure

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/kreon-core/shadow-cat-common/logc"
	"github.com/kreon-core/shadow-cat-common/utlc"
	"github.com/spf13/viper"

	"sc-player-service/infrastructure/config"
)

const (
	configPath = "configs"
	configName = "application"
	configType = "yaml"
)

func LoadEnvs(files ...string) {
	err := godotenv.Load(files...)
	if err != nil {
		logc.Warn().Strs("files", files).AnErr("warn", err).Msg("No .env files found, skipping...")
	} else {
		logc.Info().Strs("files", files).Msg(".env files loaded successfully")
	}
}

func LoadConfigs(profile ...string) (*config.Config, error) {
	cfg := &config.Config{}

	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigType(configType)

	filename := configName
	if len(profile) > 0 && !utlc.IsBlank(profile[0]) {
		filename = fmt.Sprintf("%s-%s", configName, profile[0])
	}

	v.SetConfigName(filename)

	err := v.MergeInConfig()
	if err != nil {
		logc.Warn().Str("file", fmt.Sprintf("%s.%s", filename, configType)).
			AnErr("warn", err).Msg("No configuration file found")
	} else {
		logc.Info().Str("file", fmt.Sprintf("%s.%s", filename, configType)).
			Msg("Configuration file loaded")
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	err = v.Unmarshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal_config -> %w", err)
	}

	va := validator.New()
	err = va.Struct(cfg)
	if err != nil {
		return nil, fmt.Errorf("validate_config -> %w", err)
	}

	return cfg, nil
}

// TODO: auto reload configs when the config file changes
