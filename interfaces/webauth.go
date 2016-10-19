package interfaces

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"game-tracker/models/request"
)

func (handler WebserviceHandler) Login(c *gin.Context) (string, int) {
	loginInfo := request.LoginInfo{}
	err := c.BindJSON(&loginInfo)
	if err != nil {
		return "", 400
	}

	id, err, code := handler.ProfileInteractor.FindLoginId(loginInfo.Username, loginInfo.Password)
	if err != nil {
		c.Error(err)
		return "", code
	}

	tokenString, err := createToken(id)
	if err != nil {
		c.Error(err)
		return "", 500
	}

	return tokenString, 200
}

func createToken(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": id,
	})
	tokenString, err := token.SignedString([]byte("5230"))
	return tokenString, err
}
