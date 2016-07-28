package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strconv"
)

func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("X-Auth-Key")
		if tokenString == "" {
			err := fmt.Errorf("Token cannot be empty")
			c.AbortWithError(400, err)
			return
		}

		token, err := jwt.Parse(tokenString,
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte("5230"), nil
			})
		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !(ok && token.Valid) {
			if !ok {
				err = fmt.Errorf("Error parsing claims")
			}
			if !token.Valid {
				err = fmt.Errorf("Token invalid")
			}
			c.AbortWithError(400, err)
			return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		if int(claims["id"].(float64)) != id {
			err := fmt.Errorf("Id in token and query mismatch")
			c.AbortWithError(400, err)
			return
		}
		c.Next()
	}
}
