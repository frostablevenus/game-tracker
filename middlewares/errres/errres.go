package errres

import (
	"github.com/gin-gonic/gin"
)

func ErrorHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Errors.Last() != nil {
			code := c.MustGet("code").(int)
			c.JSON(code, gin.H{
				"errors": c.Errors,
			})
			c.Abort()
		}
	}
}
