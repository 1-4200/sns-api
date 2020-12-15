package user

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
	Search(c *gin.Context)
	GetById(c *gin.Context)
	GetByIds(c *gin.Context)
}

type userHandler struct {
	l           logger.Logging
	userUseCase usecase.UserUseCase
}

func NewUserHandler(l logger.Logging, uu usecase.UserUseCase) Handler {
	return &userHandler{
		l:           l,
		userUseCase: uu,
	}
}

func (uh *userHandler) Search(c *gin.Context) {
	var q SearchForm

	q.Count, _ = strconv.Atoi(c.DefaultQuery(string(q.Count), "10"))
	q.OrderBy = c.DefaultQuery(q.OrderBy, "followers_count")

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
	users, hits, err := uh.userUseCase.Search(q.Name, q.Description, q.Language, q.FollowerMin, q.FollowerMax, q.StatusMin, q.StatusMax, q.FavoriteMin, q.FavoriteMax, q.FollowMin, q.FollowMax, q.ListMin, q.ListMax, q.SrScoreMin, q.SrScoreMax, q.StartDate, q.EndDate, q.Count, q.OrderBy)
	if err != nil {
		uh.l.Errorf(fmt.Sprintf("failed to Search: %v", err))
		c.Error(err).SetType(gin.ErrorTypePrivate).SetMeta(http.StatusNoContent)
		return
	}
	r := &Response{
		Hits: hits,
		Res:  users,
	}
	c.JSON(http.StatusOK, r)
}

func (uh *userHandler) GetById(c *gin.Context) {
	var q IDForm

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
	user, hits, err := uh.userUseCase.GetById(q.UserID, q.StartDate, q.EndDate)
	if err != nil {
		uh.l.Errorf(fmt.Sprintf("failed to GetById: %v", err))
		c.Error(err).SetType(gin.ErrorTypePrivate).SetMeta(http.StatusNoContent)
		return
	}
	r := &Response{
		Hits: hits,
		Res:  user,
	}
	c.JSON(http.StatusOK, r)
}

func (uh *userHandler) GetByIds(c *gin.Context) {
	var q IDsForm

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
	users, hits, err := uh.userUseCase.GetByIds(q.UserIDs, q.StartDate, q.EndDate)
	if err != nil {
		uh.l.Errorf(fmt.Sprintf("failed to GetByIds: %v", err))
		c.Error(err).SetType(gin.ErrorTypePrivate).SetMeta(http.StatusNoContent)
		return
	}
	r := &Response{
		Hits: hits,
		Res:  users,
	}
	c.JSON(http.StatusOK, r)
}
