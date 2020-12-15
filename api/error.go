package api

import (
	"github.com/gin-gonic/gin"
)

func (s *server) HandleError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if err := c.Errors.ByType(gin.ErrorTypePrivate).Last(); err != nil {
			statusCode := c.Errors.ByType(gin.ErrorTypePrivate).Last().Meta.(int)
			c.AbortWithStatusJSON(statusCode, gin.H{
				"error": err.Error(),
			})
		}
	}
}
