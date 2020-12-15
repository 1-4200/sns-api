package hashtag

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"sns-api/logger"
	"sns-api/usecase"
	"strconv"
)

type Handler interface {
	Get(c *gin.Context)
	Search(c *gin.Context)
}

type hashtagHandler struct {
	l              logger.Logging
	hashtagUseCase usecase.HashtagUseCase
}

func NewHashtagHandler(l logger.Logging, hu usecase.HashtagUseCase) Handler {
	return &hashtagHandler{
		l:              l,
		hashtagUseCase: hu,
	}
}

func (hh *hashtagHandler) Get(c *gin.Context) {
	var q Form
	q.Count, _ = strconv.Atoi(c.DefaultQuery(string(q.Count), "1000"))
	q.RetweetMin, _ = strconv.Atoi(c.DefaultQuery(string(q.RetweetMin), "0"))
	q.QuoteMin, _ = strconv.Atoi(c.DefaultQuery(string(q.QuoteMin), "0"))
	q.FavoriteMin, _ = strconv.Atoi(c.DefaultQuery(string(q.FavoriteMin), "0"))

	if err := c.ShouldBind(&q); err != nil {
		switch e := err.(type) {
		case validator.FieldError:
			for _, fieldErr := range err.(validator.ValidationErrors) {
				c.Error(errors.New(fmt.Sprint(fieldErr))).SetType(gin.ErrorTypePrivate).SetMeta(http.StatusBadRequest)
				return
			}
		default:
			c.Error(e).SetType(gin.ErrorTypePrivate).SetMeta(http.StatusBadRequest)
			return
		}
	}
	hashtags, hits, err := hh.hashtagUseCase.Get(q.Keyword, q.Hashtag, q.TweetType, q.RetweetMin, q.RetweetMax, q.QuoteMin, q.QuoteMax, q.FavoriteMin, q.FavoriteMax, q.UserInclude, q.UserExclude, q.HashtagInclude, q.HashtagExclude, q.UserFollowerMin, q.UserFollowerMax, q.UserStatusMin, q.UserStatusMax, q.Count, q.StartDate, q.EndDate)
	if err != nil {
		hh.l.Errorf(fmt.Sprintf("failed to Get: %v", err))
		c.Error(err).SetType(gin.ErrorTypePrivate).SetMeta(http.StatusNoContent)
		return
	}
	r := &Response{
		Hits: hits,
		Res:  hashtags,
	}
	c.JSON(http.StatusOK, r)
}

func (hh *hashtagHandler) Search(c *gin.Context) {
	var q SearchForm
	q.Count, _ = strconv.Atoi(c.DefaultQuery(string(q.Count), "10"))

	if err := c.ShouldBind(&q); err != nil {
		switch e := err.(type) {
		case validator.FieldError:
			for _, fieldErr := range err.(validator.ValidationErrors) {
				c.Error(errors.New(fmt.Sprint(fieldErr))).SetType(gin.ErrorTypePrivate).SetMeta(http.StatusBadRequest)
				return
			}
		default:
			c.Error(e).SetType(gin.ErrorTypePrivate).SetMeta(http.StatusBadRequest)
			return
		}
	}
	hashtags, hits, err := hh.hashtagUseCase.Search(q.Hashtag, q.StartDate, q.EndDate, q.Count)
	if err != nil {
		hh.l.Errorf(fmt.Sprintf("failed to Search: %v", err))
		c.Error(err).SetType(gin.ErrorTypePrivate).SetMeta(http.StatusNoContent)
		return
	}
	r := &Response{
		Hits: hits,
		Res:  hashtags,
	}
	c.JSON(http.StatusOK, r)
}
