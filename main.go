package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"sns-api/api"
	"sns-api/config"
	"sns-api/logger"
)

var AppEnvironment string

func init()  {
	AppEnvironment = os.Getenv("APP_ENVIRONMENT")
}

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(-1)
	}
}

func setup() (*gin.Engine, *config.Config) {
	c, err := config.NewConfig(AppEnvironment)
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}
	l, err := logger.NewLogger(c)
	if err != nil {
		log.Fatalf("cannot create logger instance: %v", err)
	}
	if c.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())

	server := api.NewServer(r, c, l)
	server.NewRouter()
	return r, c
}

func run() error {
	r, c := setup()
	return r.Run(fmt.Sprintf(":%s", c.Port))
}
