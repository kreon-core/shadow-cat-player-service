package i12e

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	tul "github.com/kreon-core/shadow-cat-common"
	"github.com/kreon-core/shadow-cat-common/logc"
	"github.com/spf13/viper"

	"sc-player-service/i12e/config"
)

const (
	configPath = "configs"
	configName = "application"
	configType = "yaml"
)

func LoadEnvs(files ...string) {
	_ = godotenv.Load()
	_ = godotenv.Load(files...)
}

func LoadConfigs(profile ...string) (*config.Config, error) {
	cfg := &config.Config{}

	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigType(configType)

	filename := configName
	if len(profile) > 0 && !tul.IsBlank(profile[0]) {
		filename = fmt.Sprintf("%s-%s", configName, profile[0])
	}

	v.SetConfigName(filename)

	err := v.MergeInConfig()
	if err != nil {
		logc.Warn("No configuration file found", "file", fmt.Sprintf("%s.%s", filename, configType), err)
	} else {
		logc.Info("Loaded configuration", "file", fmt.Sprintf("%s.%s", filename, configType))
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
