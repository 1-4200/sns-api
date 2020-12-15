package tweet

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"sns-api/handler"
	"sns-api/logger"
	"sns-api/usecase"
	"strconv"
)

type Handler interface {
	Get(c *gin.Context)
	GetByUser(c *gin.Context)
	GetByUsers(c *gin.Context)
	GetByDomain(c *gin.Context)
	GetByMediaType(c *gin.Context)
	GetTransitionByUser(c *gin.Context)
}

type tweetHandler struct {
	l            logger.Logging
	tweetUseCase usecase.TweetUseCase
}

func NewTweetHandler(l logger.Logging, tu usecase.TweetUseCase) Handler {
	return &tweetHandler{
		l:            l,
		tweetUseCase: tu,
	}
}

func (th *tweetHandler) Get(c *gin.Context) {
	tweets, err := th.tweetUseCase.Get()
	if err != nil {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, tweets)
}

func (th *tweetHandler) GetByUser(c *gin.Context) {
	var q UserForm

	q.Count, _ = strconv.Atoi(c.DefaultQuery(string(q.Count), "1"))
	q.OrderBy = c.DefaultQuery(q.OrderBy, "favorite_count")

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

	tweets, hits, err := th.tweetUseCase.GetByUser(q.UserID, handler.ConvertTime(q.StartDate), handler.ConvertTime(q.EndDate), q.Count, q.OrderBy)
	if err != nil {
		th.l.Errorf(fmt.Sprintf("failed to GetByUser: %v", err))
		c.Error(err).SetType(gin.ErrorTypePrivate).SetMeta(http.StatusNoContent)
		return
	}
	r := &Response{
		Hits: hits,
		Res:  tweets,
	}
	c.JSON(http.StatusOK, r)
}

func (th *tweetHandler) GetByUsers(c *gin.Context) {
	var q UsersForm

	q.Count, _ = strconv.Atoi(c.DefaultQuery(string(q.Count), "1"))
	q.OrderBy = c.DefaultQuery(q.OrderBy, "created_at")

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
	tweets, hits, err := th.tweetUseCase.GetByUsers(q.UserIDs, q.StartDate, q.EndDate, q.Count, q.OrderBy)
	if err != nil {
		th.l.Errorf(fmt.Sprintf("failed to GetByUsers: %v", err))
		c.Error(err).SetType(gin.ErrorTypePrivate).SetMeta(http.StatusNoContent)
		return
	}
	r := &Response{
		Hits: hits,
		Res:  tweets,
	}
	c.JSON(http.StatusOK, r)
}

func (th *tweetHandler) GetByDomain(c *gin.Context) {
	var q URLForm

	q.Count, _ = strconv.Atoi(c.DefaultQuery(string(q.Count), "1"))
	q.OrderBy = c.DefaultQuery(q.OrderBy, "favorite_count")

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

	tweets, hits, urlInfo, err := th.tweetUseCase.GetByDomain(q.UserID, handler.ConvertTime(q.StartDate), handler.ConvertTime(q.EndDate), q.Count, q.OrderBy, q.Domain)
	if err != nil {
		th.l.Errorf(fmt.Sprintf("failed to GetByDomain: %v", err))
		c.Error(err).SetType(gin.ErrorTypePrivate).SetMeta(http.StatusNoContent)
		return
	}
	r := &ResponseDomain{
		Hits:    hits,
		Tweets:  tweets,
		UrlInfo: urlInfo,
	}
	c.JSON(http.StatusOK, r)
	th.l.Info("function handler.GetByDomain done")
}

func (th *tweetHandler) GetByMediaType(c *gin.Context) {
	var q MediaForm

	q.Count, _ = strconv.Atoi(c.DefaultQuery(string(q.Count), "1"))
	q.OrderBy = c.DefaultQuery(q.OrderBy, "favorite_count")

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

	tweets, hits, media, err := th.tweetUseCase.GetByMediaType(q.UserID, handler.ConvertTime(q.StartDate), handler.ConvertTime(q.EndDate), q.Count, q.OrderBy, q.MediaType)
	if err != nil {
		th.l.Errorf(fmt.Sprintf("failed to GetByMediaType: %v", err))
		c.Error(err).SetType(gin.ErrorTypePrivate).SetMeta(http.StatusNoContent)
		return
	}
	r := &ResponseMedia{
		Hits:   hits,
		Tweets: tweets,
		Media:  media,
	}
	c.JSON(http.StatusOK, r)
	th.l.Info("function handler.GetByMedia done")
}

func (th *tweetHandler) GetTransitionByUser(c *gin.Context) {
	var q TransitionForm

	q.Count, _ = strconv.Atoi(c.DefaultQuery(string(q.Count), "100"))

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
	transitions, err := th.tweetUseCase.GetTransitionByUser(q.UserID, handler.ConvertDate(q.StartDate), handler.ConvertDate(q.EndDate), q.Count)
	if err != nil {
		th.l.Errorf(fmt.Sprintf("failed to GetTransitionByUser: %v", err))
		c.Error(err).SetType(gin.ErrorTypePrivate).SetMeta(http.StatusNoContent)
		return
	}
	r := &Response{
		Hits: len(transitions),
		Res:  transitions,
	}
	c.JSON(http.StatusOK, r)
}
