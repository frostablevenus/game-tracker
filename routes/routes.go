package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"

	"game-tracker/interfaces"
)

func CreateEngine(webserviceHandler interfaces.WebserviceHandler) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	engine.POST("/add_user", func(c *gin.Context) {
		err, code := webserviceHandler.AddUser(c)
		if err != nil {
			fmt.Println(err)
		}
		c.Status(code)
	})

	users := engine.Group("/users/:id")
	users.GET("", func(c *gin.Context) {
		err, code := webserviceHandler.ShowUser(c)
		if err != nil {
			fmt.Println(err)
		}
		c.Status(code)
	})
	users.DELETE("", func(c *gin.Context) {
		err, code := webserviceHandler.RemoveUser(c)
		if err != nil {
			fmt.Println(err)
		}
		c.Status(code)
	})
	users.GET("/info", func(c *gin.Context) {
		err, code := webserviceHandler.ShowUserInfo(c)
		if err != nil {
			fmt.Println(err)
		}
		c.Status(code)
	})
	users.PUT("/info", func(c *gin.Context) {
		err, code := webserviceHandler.EditUserInfo(c)
		if err != nil {
			fmt.Println(err)
		}
		c.Status(code)
	})
	users.POST("/add_library", func(c *gin.Context) {
		err, code := webserviceHandler.AddLibrary(c)
		if err != nil {
			fmt.Println(err)
		}
		c.Status(code)
	})

	libraries := users.Group("/libraries/:libId")
	libraries.GET("", func(c *gin.Context) {
		err, code := webserviceHandler.ShowLibrary(c)
		if err != nil {
			fmt.Println(err)
		}
		c.Status(code)
	})
	libraries.DELETE("", func(c *gin.Context) {
		err, code := webserviceHandler.RemoveLibrary(c)
		if err != nil {
			fmt.Println(err)
		}
		c.Status(code)
	})
	libraries.POST("/add_game", func(c *gin.Context) {
		err, code := webserviceHandler.AddGame(c)
		if err != nil {
			fmt.Println(err)
		}
		c.Status(code)
	})

	games := libraries.Group("/games/:gameId")
	games.DELETE("", func(c *gin.Context) {
		err, code := webserviceHandler.RemoveGame(c)
		if err != nil {
			fmt.Println(err)
		}
		c.Status(code)
	})
	return engine
}
