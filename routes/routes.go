package routes

import (
	"github.com/gin-gonic/gin"

	"game-tracker/interfaces"
)

func CreateEngine(webserviceHandler interfaces.WebserviceHandler) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	engine.POST("/add_user", func(c *gin.Context) {
		err, code, message := webserviceHandler.AddUser(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.String(code, message)
		}
	})

	users := engine.Group("/users/:id")
	users.GET("", func(c *gin.Context) {
		err, code, message := webserviceHandler.ShowUser(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.JSON(code, message)
		}
	})
	users.DELETE("", func(c *gin.Context) {
		err, code, message := webserviceHandler.RemoveUser(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.String(code, message)
		}
	})
	users.GET("/info", func(c *gin.Context) {
		err, code, message := webserviceHandler.ShowUserInfo(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.JSON(code, message)
		}
	})
	users.PUT("/info", func(c *gin.Context) {
		err, code, message := webserviceHandler.EditUserInfo(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.String(code, message)
		}
	})
	users.POST("/add_library", func(c *gin.Context) {
		err, code, message := webserviceHandler.AddLibrary(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.String(code, message)
		}
	})

	libraries := users.Group("/libraries/:libId")
	libraries.GET("", func(c *gin.Context) {
		err, code, message := webserviceHandler.ShowLibrary(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.JSON(code, message)
		}
	})
	libraries.DELETE("", func(c *gin.Context) {
		err, code, message := webserviceHandler.RemoveLibrary(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.JSON(code, message)
		}
	})
	libraries.POST("/add_game", func(c *gin.Context) {
		err, code, message := webserviceHandler.AddGame(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.JSON(code, message)
		}
	})

	games := libraries.Group("/games/:gameId")
	games.DELETE("", func(c *gin.Context) {
		err, code, message := webserviceHandler.RemoveGame(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.JSON(code, message)
		}
	})
	return engine
}
