package api

import (
	"database/sql"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
	"sns-api/config"
	"sns-api/logger"
	"strings"
	"time"
)

type server struct {
	router *gin.Engine
	config *config.Config
	logger logger.Logging
	es     *elasticsearch.Client
	corpus *sql.DB
}

func NewServer(e *gin.Engine, c *config.Config, l logger.Logging) *server {
	return &server{
		router: e,
		config: c,
		logger: l,
	}
}

func (s *server) HandleAccessLog() gin.HandlerFunc {
	f, err := os.OpenFile("log/access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		s.logger.Fatalf(fmt.Sprintf("error opening file: %v", err))
	}
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				strings.Replace(param.ErrorMessage, "\n", "", -1),
			)
		},
		Output:    f,
		SkipPaths: nil,
	})
}

func (s *server) NewElasticSearchClient() gin.HandlerFunc {
	cfg := elasticsearch.Config{
		Addresses: []string{s.config.DB.ElasticSearch.Address},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		s.logger.Fatalf(fmt.Sprintf("cannot create elasticserch client: %v", err))
	}
	s.es = es
	return func(c *gin.Context) {
		_, err := s.es.Ping()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Next()
	}
}

func (s *server) NewCorpusDatabaseClient() gin.HandlerFunc {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", s.config.DB.Corpus.Username, s.config.DB.Corpus.Password, s.config.DB.Corpus.Host, s.config.DB.Corpus.Port, s.config.DB.Corpus.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		s.logger.Fatalf(fmt.Sprintf("cannot create corpus client: %v", err))
	}
	s.corpus = db
	// no need to close corpus db connection here
	return func(c *gin.Context) {
		if err = s.corpus.Ping(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Next()
	}
}
