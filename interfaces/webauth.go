package interfaces

import (
	//"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"game-tracker/models/request"
	"game-tracker/models/result"
)

func (handler WebserviceHandler) Login(c *gin.Context) (string, result.Errors, int) {
	loginInfo := request.LoginInfo{}
	errors := result.Errors{}
	errMsg := c.BindJSON(&loginInfo)
	if errMsg != nil {
		err := result.Error{Message: errMsg}
		errors.Errs = append(errors.Errs, err)
		return "", errors, 400
	}

	id, errMsg, code := handler.ProfileInteractor.FindLoginId(loginInfo.Username, loginInfo.Password)
	if errMsg != nil {
		err := result.Error{Message: errMsg}
		errors.Errs = append(errors.Errs, err)
		return "", errors, code
	}

	tokenString, errMsg := createToken(id)
	if errMsg != nil {
		err := result.Error{Message: errMsg}
		errors.Errs = append(errors.Errs, err)
		return "", errors, 500
	}

	return tokenString, errors, 200
}

func createToken(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": id,
	})
	tokenString, err := token.SignedString([]byte("5230"))
	return tokenString, err
}
