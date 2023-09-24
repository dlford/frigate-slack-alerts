package config

import (
	"app/sharedTypes"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

func GetConfig() (*sharedTypes.Config, error) {
	godotenv.Load()

	var cfg sharedTypes.Config
	err := envconfig.Process("fsa", &cfg)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(cfg.Frigate_Internal_BaseURL, "/") {
		strings.TrimSuffix(cfg.Frigate_Internal_BaseURL, "/")
	}

	if cfg.Filter_Config_File != "" {
		file, err := os.ReadFile(cfg.Filter_Config_File)
		if err != nil {
			panic(err)
		}

		err = yaml.Unmarshal(file, &cfg.Filters)
		if err != nil {
			panic(err)
		}
	}

	return &cfg, nil
}
