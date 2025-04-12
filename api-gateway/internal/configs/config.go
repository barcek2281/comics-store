package configs

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Port int `yaml:"port"`
	StoragePath string `yaml:"storage_path"`
}

func MustLoad(configPath string) *Config {
	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		panic(err)
	}

	return &config
}