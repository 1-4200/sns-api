package config

import (
	"github.com/jinzhu/configor"
	"strings"
)

type Config struct {
	AppName string `default:"sns-api"`
	Port    string `default:"8080"`
	Env     string `default:"dev"`
	Logger  struct {
		Use         string `default:"zap"`
		Environment string `default:"dev"`
		LogLevel    string `default:"debug"`
		FileName    string `default:"log/app.log"`
	}
	DB struct {
		Corpus struct{
			Host     string `default:"mysql"`
			Port     string `default:"3306"`
			Username string `default:"mysql"`
			Password string `default:"mysql"`
			Database string `default:"app"`
		}
		ElasticSearch struct {
				Address string `default:"http://localhost:8443"`
		}
	}
}

func NewConfig(env string) (*Config, error) {
	var err error
	c := &Config{}
	switch strings.ToLower(env) {
	case "dev", "development":
		err = configor.Load(c, "config/config.dev.yml")
	case "prod", "production":
		err = configor.Load(c, "config/config.prod.yml")
	default:
		err = configor.Load(c, "config/config.dev.yml")
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}
