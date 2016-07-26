package routes

import (
	"github.com/gin-gonic/gin"

	"game-tracker/interfaces"
	"game-tracker/middlewares/auth"
)

func CreateEngine(webserviceHandler interfaces.WebserviceHandler) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	engine.POST("/login", func(c *gin.Context) {
		tokenString, err, code := webserviceHandler.Login(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.JSON(code, gin.H{"token": tokenString})
		}
	})

	unAuth := engine.Group("/users")
	unAuth.GET("/:id", func(c *gin.Context) {
		err, code, message := webserviceHandler.ShowUser(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.JSON(code, message)
		}
	})
	unAuth.POST("", func(c *gin.Context) {
		err, code, message := webserviceHandler.AddUser(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.String(code, message)
		}
	})
	unAuth.GET("/:id/info", func(c *gin.Context) {
		err, code, message := webserviceHandler.ShowUserInfo(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.JSON(code, message)
		}
	})

	authorized := engine.Group("/users/:id")
	authorized.Use(auth.CheckToken())

	users := authorized.Group("")
	users.DELETE("", func(c *gin.Context) {
		err, code, message := webserviceHandler.RemoveUser(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.String(code, message)
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

	libraries := users.Group("/libraries")
	libraries.GET("/:libId", func(c *gin.Context) {
		err, code, message := webserviceHandler.ShowLibrary(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.JSON(code, message)
		}
	})
	libraries.POST("", func(c *gin.Context) {
		err, code, message := webserviceHandler.AddLibrary(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.String(code, message)
		}
	})
	libraries.DELETE("/:libId", func(c *gin.Context) {
		err, code, message := webserviceHandler.RemoveLibrary(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.JSON(code, message)
		}
	})

	games := libraries.Group("/:libId/games")
	games.POST("", func(c *gin.Context) {
		err, code, message := webserviceHandler.AddGame(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.String(code, message)
		}
	})
	games.POST("/:gameId", func(c *gin.Context) {
		err, code, message := webserviceHandler.PickGame(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.String(code, message)
		}
	})
	games.DELETE("/:gameId", func(c *gin.Context) {
		err, code, message := webserviceHandler.RemoveGame(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.JSON(code, message)
		}
	})
	return engine
}
