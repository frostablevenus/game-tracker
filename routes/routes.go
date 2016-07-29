package routes

import (
	"github.com/gin-gonic/gin"

	"game-tracker/interfaces"
	"game-tracker/middlewares/auth"
	"game-tracker/middlewares/errres"
	res "game-tracker/models/responses"
)

func CreateEngine(webserviceHandler interfaces.WebserviceHandler) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	engine.Use(errres.ErrorHandle())

	engine.POST("/login", func(c *gin.Context) {
		tokenString, code := webserviceHandler.Login(c)
		c.Set("code", code)
		if c.Errors == nil {
			token := res.ViewToken(tokenString)
			c.JSON(201, token)
		}
	})

	unAuth := engine.Group("/users")
	unAuth.GET("/:id", func(c *gin.Context) {
		code, message := webserviceHandler.ShowUser(c)
		c.Set("code", code)
		if c.Errors == nil {
			libraries := res.ViewLibraries(message.LibraryIds)
			users := res.ViewUser(message.Id, message.Name, libraries)
			c.JSON(200, users)
		}
	})
	unAuth.POST("", func(c *gin.Context) {
		code, message := webserviceHandler.AddUser(c)
		c.Set("code", code)
		if c.Errors == nil {
			users := res.ViewUser(message.Id, message.Name, nil)
			c.JSON(201, users)
		}
	})
	unAuth.GET("/:id/info", func(c *gin.Context) {
		code, message := webserviceHandler.ShowUserInfo(c)
		c.Set("code", code)
		if c.Errors == nil {
			info := res.ViewInfo(message.Info, message.Id)
			c.JSON(200, info)
		}
	})

	authorized := engine.Group("/users/:id")
	authorized.Use(auth.CheckToken())

	users := authorized.Group("")
	users.DELETE("", func(c *gin.Context) {
		code, _ := webserviceHandler.RemoveUser(c)
		c.Set("code", code)
		if c.Errors == nil {
			c.Status(204)
		}
	})
	users.PUT("/info", func(c *gin.Context) {
		code, message := webserviceHandler.EditUserInfo(c)
		c.Set("code", code)
		if c.Errors == nil {
			info := res.ViewInfo(message.Info, message.Id)
			c.JSON(201, info)
		}
	})

	libraries := users.Group("/libraries")
	libraries.GET("/:libId", func(c *gin.Context) {
		code, message := webserviceHandler.ShowLibrary(c)
		c.Set("code", code)
		if c.Errors == nil {
			games := res.ViewGames(message.GamesIds)
			library := res.ViewLibrary(message.UserId, message.Id, games)
			c.JSON(200, library)
		}
	})
	libraries.POST("", func(c *gin.Context) {
		code, message := webserviceHandler.AddLibrary(c)
		c.Set("code", code)
		if c.Errors == nil {
			library := res.ViewLibrary(message.UserId, message.Id, nil)
			c.JSON(201, library)
		}
	})
	libraries.DELETE("/:libId", func(c *gin.Context) {
		code, _ := webserviceHandler.RemoveLibrary(c)
		c.Set("code", code)
		if c.Errors == nil {
			c.Status(204)
		}
	})

	games := libraries.Group("/:libId/games")
	games.POST("", func(c *gin.Context) {
		code, message := webserviceHandler.AddGame(c)
		c.Set("code", code)
		if c.Errors == nil {
			game := res.ViewGame(message.UserId, message.LibraryId, message.Id,
				message.Name, message.Producer, message.Value)
			c.JSON(code, game)
		}
	})
	games.POST("/:gameId", func(c *gin.Context) {
		code, message := webserviceHandler.PickGame(c)
		c.Set("code", code)
		if c.Errors == nil {
			game := res.ViewGame(message.UserId, message.LibraryId, message.Id, "", "", 0)
			c.JSON(code, game)
		}
	})
	games.DELETE("/:gameId", func(c *gin.Context) {
		code, _ := webserviceHandler.RemoveGame(c)
		c.Set("code", code)
		if c.Errors == nil {
			c.Status(204)
		}
	})
	return engine
}
