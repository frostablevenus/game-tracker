package auth

import (
	"fmt"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("X-Auth-Key")
		if tokenString == "" {
			err := fmt.Errorf("Token empty")
			c.AbortWithError(400, err)
		}
		token, err := jwt.Parse(tokenString,
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte("5230"), nil
			})

		claims, ok := token.Claims.(jwt.MapClaims)
		if !(ok && token.Valid) {
			c.AbortWithError(403, err)
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithError(400, err)
		}
		if int(claims["id"].(float64)) != id {
			err := fmt.Errorf("Id in token and query mismatch")
			c.AbortWithError(403, err)
		}
		c.Next()
	}
}
