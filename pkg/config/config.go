package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	HTTPConfig    `yaml:"http"`
	LogConfig     `yaml:"log"`
	DBConfig      `yaml:"db"`
	WeatherConfig `yaml:"weather"`
}

type HTTPConfig struct {
	Port   int    `yaml:"port"`
	Secret string `yaml:"secret"`
}

type LogConfig struct {
	Level    string `yaml:"level"`
	Encoding string `yaml:"encoding"`
}

type DBConfig struct {
	ConnURL string `yaml:"databaseURL"`
}

type WeatherConfig struct {
	Mock     bool   `yaml:"mock"`
	APIKey   string `yaml:"apiKey"`
	FileName string `yaml:"file"`
}

func LoadConfig(fileName string) (*Config, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read config file in %v", fileName)
	}

	cfg := Config{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, errors.WithStack(err)
	}

	return &cfg, nil
}
