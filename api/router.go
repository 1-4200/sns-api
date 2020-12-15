package api

import (
	"github.com/gin-gonic/gin"
	"sns-api/handler/hashtag"
	"sns-api/handler/tweet"
	"sns-api/handler/user"
	"sns-api/infrastructure/elastic"
	"sns-api/infrastructure/mysql/corpus"
	"sns-api/usecase"
)

func (s *server) NewRouter() {
	s.router.Use(s.HandleAccessLog())
	s.router.Use(s.HandleError())
	s.router.Use(s.NewElasticSearchClient())
	s.router.Use(s.NewCorpusDatabaseClient())

	apiV1 := s.router.Group("api/v1")
	s.healthRoutes(apiV1)
	s.tweetsRoutes(apiV1)
	s.hashtagsRoutes(apiV1)
	s.usersRoutes(apiV1)
}

func (s *server) healthRoutes(api *gin.RouterGroup) {
	healthRoutes := api.Group("/health")
	{
		healthRoutes.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "success",
			})
		})
	}
}

func (s *server) tweetsRoutes(api *gin.RouterGroup) {
	tweetsRoutes := api.Group("/tweets")
	{
		tweetRepository := elastic.NewTweetRepository(s.logger, s.es)
		corpusTweetRepository := corpus.NewTweetRepository(s.logger, s.corpus)
		tweetUseCase := usecase.NewTweetUseCase(s.logger, tweetRepository, corpusTweetRepository)
		tweetHandler := tweet.NewTweetHandler(s.logger, tweetUseCase)

		tweetsRoutes.GET("/", tweetHandler.Get)
		tweetsRoutes.GET("/user", tweetHandler.GetByUser)
		tweetsRoutes.GET("/users", tweetHandler.GetByUsers)
		tweetsRoutes.POST("/users", tweetHandler.GetByUsers)
		tweetsRoutes.GET("/domain", tweetHandler.GetByDomain)
		tweetsRoutes.GET("/media", tweetHandler.GetByMediaType)
		tweetsRoutes.GET("/transition", tweetHandler.GetTransitionByUser)
	}
}

func (s *server) hashtagsRoutes(api *gin.RouterGroup) {
	hashtagsRoutes := api.Group("/hashtags")
	{
		hashtagRepository := elastic.NewHashtagRepository(s.logger, s.es)
		hashtagUseCase := usecase.NewHashtagUseCase(s.logger, hashtagRepository)
		hashtagHandler := hashtag.NewHashtagHandler(s.logger, hashtagUseCase)

		hashtagsRoutes.GET("/", hashtagHandler.Get)
		hashtagsRoutes.GET("/search", hashtagHandler.Search)
	}
}

func (s *server) usersRoutes(api *gin.RouterGroup) {
	usersRoutes := api.Group("/users")
	{
		userRepository := elastic.NewUserRepository(s.logger, s.es)
		userUseCase := usecase.NewUserUseCase(s.logger, userRepository)
		userHandler := user.NewUserHandler(s.logger, userUseCase)

		usersRoutes.GET("/search", userHandler.Search)
		usersRoutes.GET("/id", userHandler.GetById)
		usersRoutes.GET("/ids", userHandler.GetByIds)
		usersRoutes.POST("/ids", userHandler.GetByIds)
	}
}
