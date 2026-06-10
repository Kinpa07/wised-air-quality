package app

import (
	"context"
	"log"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfigFromFile(cfgFilePath string, cfg *Config) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(cfgFilePath)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(cfg); err != nil {
		return err
	}
	return nil
}

func LoadConfig(ctx context.Context, cfgFilePath string, cfg *Config) error {
	log.Printf("Loading config from path: %s", cfgFilePath)
	if err := LoadConfigFromFile(cfgFilePath, cfg); err != nil {
		return err
	}
	return nil
}
