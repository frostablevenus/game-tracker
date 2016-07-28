package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strconv"

	"game-tracker/models/result"
)

func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("X-Auth-Key")
		if tokenString == "" {
			errMsg := fmt.Errorf("Token cannot be empty")
			err := result.Error{Message: errMsg}
			c.JSON(400, gin.H{
				"errors": []gin.H{
					gin.H{
						"message": err.Message.Error(),
					},
				},
			})
			c.Abort()
			return
		}

		token, errMsg := jwt.Parse(tokenString,
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte("5230"), nil
			})
		if errMsg != nil {
			err := result.Error{Message: errMsg}
			c.JSON(400, gin.H{
				"errors": []gin.H{
					gin.H{
						"message": err.Message.Error(),
					},
				},
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !(ok && token.Valid) {
			if !ok {
				errMsg = fmt.Errorf("Error parsing claims")
			}
			if !token.Valid {
				errMsg = fmt.Errorf("Token invalid")
			}
			err := result.Error{Message: errMsg}
			c.JSON(400, gin.H{
				"errors": []gin.H{
					gin.H{
						"message": err.Message.Error(),
					},
				},
			})
			c.Abort()
			return
		}

		id, errMsg := strconv.Atoi(c.Param("id"))
		if errMsg != nil {
			err := result.Error{Message: errMsg}
			c.JSON(400, gin.H{
				"errors": []gin.H{
					gin.H{
						"message": err.Message.Error(),
					},
				},
			})
			c.Abort()
			return
		}

		if int(claims["id"].(float64)) != id {
			errMsg := fmt.Errorf("Id in token and query mismatch")
			err := result.Error{Message: errMsg}
			c.JSON(400, gin.H{
				"errors": []gin.H{
					gin.H{
						"message": err.Message.Error(),
					},
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
