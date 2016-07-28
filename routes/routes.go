package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"

	"game-tracker/interfaces"
	"game-tracker/middlewares/auth"
	"game-tracker/middlewares/errres"
)

func CreateEngine(webserviceHandler interfaces.WebserviceHandler) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	engine.Use(errres.ErrorHandle())

	engine.POST("/login", func(c *gin.Context) {
		tokenString, code := webserviceHandler.Login(c)
		c.Set("code", code)
		if c.Errors == nil {
			c.JSON(201, gin.H{
				"data": gin.H{
					"type": "token",
					"attributes": gin.H{
						"tokenString": tokenString,
					},
				},
			})
		}
	})

	unAuth := engine.Group("/users")
	unAuth.GET("/:id", func(c *gin.Context) {
		code, message := webserviceHandler.ShowUser(c)
		c.Set("code", code)
		if c.Errors == nil {
			libraries := []gin.H{}
			for _, libraryId := range message.LibraryIds {
				libraries = append(libraries, gin.H{
					"library": gin.H{
						"data": gin.H{
							"type": "libraries",
							"id":   libraryId,
						},
					},
				})
			}
			c.JSON(200, gin.H{
				"links": gin.H{
					"self": fmt.Sprintf("http://localhost:8080/users/%d", message.Id),
				},
				"data": gin.H{
					"type": "users",
					"id":   message.Id,
					"attributes": gin.H{
						"name": message.Name,
					},
					"relationships": gin.H{
						"libraries": libraries,
					},
				},
			})
		}
	})
	unAuth.POST("", func(c *gin.Context) {
		code, message := webserviceHandler.AddUser(c)
		c.Set("code", code)
		if c.Errors == nil {
			c.JSON(201, gin.H{
				"data": gin.H{
					"type": "users",
					"id":   message.Id,
					"attributes": gin.H{
						"name": message.Name,
					},
					"links": gin.H{
						"self": fmt.Sprintf("http://localhost:8080/users/%d", message.Id),
					},
				},
			})
		}
	})
	unAuth.GET("/:id/info", func(c *gin.Context) {
		code, message := webserviceHandler.ShowUserInfo(c)
		c.Set("code", code)
		if c.Errors == nil {
			c.JSON(200, gin.H{
				"links": gin.H{
					"self":    fmt.Sprintf("http://localhost:8080/users/%d/info", message.Id),
					"related": fmt.Sprintf("http://localhost:8080/users/%d", message.Id),
				},
				"data": gin.H{
					"type": "info",
					"id":   message.Id,
					"attributes": gin.H{
						"content": message.Info,
					},
					"relationships": gin.H{
						"owner": gin.H{
							"data": gin.H{
								"type": "users",
								"id":   message.Id,
							},
						},
					},
				},
			})
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
			c.JSON(201, gin.H{
				"data": gin.H{
					"type": "info",
					"id":   message.Id,
					"attributes": gin.H{
						"content": message.Info,
					},
					"relationships": gin.H{
						"owner": gin.H{
							"data": gin.H{
								"type": "users",
								"id":   message.Id,
							},
						},
					},
					"links": gin.H{
						"self":    fmt.Sprintf("http://localhost:8080/users/%d/info", message.Id),
						"related": fmt.Sprintf("http://localhost:8080/users/%d", message.Id),
					},
				},
			})
		}
	})

	libraries := users.Group("/libraries")
	libraries.GET("/:libId", func(c *gin.Context) {
		code, message := webserviceHandler.ShowLibrary(c)
		c.Set("code", code)
		if c.Errors == nil {
			games := []gin.H{}
			for _, gameId := range message.GamesIds {
				games = append(games, gin.H{
					"game": gin.H{
						"data": gin.H{
							"type": "games",
							"id":   gameId,
						},
					},
				})
			}
			c.JSON(200, gin.H{
				"links": gin.H{
					"self": fmt.Sprintf("http://localhost:8080/users/%d/libraries/%d",
						message.UserId, message.Id),
					"related": fmt.Sprintf("http://localhost:8080/users/%d",
						message.UserId),
				},
				"data": gin.H{
					"type": "libraries",
					"id":   message.Id,
					"relationships": gin.H{
						"games": games,
						"owner": gin.H{
							"data": gin.H{
								"type": "users",
								"id":   message.UserId,
							},
						},
					},
				},
			})
		}
	})
	libraries.POST("", func(c *gin.Context) {
		code, message := webserviceHandler.AddLibrary(c)
		c.Set("code", code)
		if c.Errors == nil {
			c.JSON(201, gin.H{
				"data": gin.H{
					"type": "libraries",
					"id":   message.Id,
					"relationships": gin.H{
						"owner": gin.H{
							"data": gin.H{
								"type": "users",
								"id":   message.UserId,
							},
						},
					},
					"links": gin.H{
						"self": fmt.Sprintf("http://localhost:8080/users/%d/libraries/%d",
							message.UserId, message.Id),
						"related": fmt.Sprintf("http://localhost:8080/users/%d", message.UserId),
					},
				},
			})
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
			c.JSON(code, gin.H{
				"data": gin.H{
					"type": "games",
					"id":   message.Id,
					"attributes": gin.H{
						"name":     message.Name,
						"producer": message.Producer,
						"value":    message.Value,
					},
					"relationships": gin.H{
						"library": gin.H{
							"data": gin.H{
								"type": "libraries",
								"id":   message.LibraryId,
							},
						},
					},
					"links": gin.H{
						"self": fmt.Sprintf("http://localhost:8080/users/%d/libraries/%d/games/%d",
							message.UserId, message.LibraryId, message.Id),
						"related": fmt.Sprintf("http://localhost:8080/users/%d/libraries/%d",
							message.UserId, message.LibraryId),
					},
				},
			})
		}
	})
	games.POST("/:gameId", func(c *gin.Context) {
		code, message := webserviceHandler.PickGame(c)
		c.Set("code", code)
		if c.Errors == nil {
			c.JSON(code, gin.H{
				"data": gin.H{
					"type": "games",
					"id":   message.Id,
					"relationships": gin.H{
						"library": gin.H{
							"data": gin.H{
								"type": "libraries",
								"id":   message.LibraryId,
							},
						},
					},
					"links": gin.H{
						"self": fmt.Sprintf("http://localhost:8080/users/%d/libraries/%d/games/%d",
							message.UserId, message.LibraryId, message.Id),
						"related": fmt.Sprintf("http://localhost:8080/users/%d/libraries/%d",
							message.UserId, message.LibraryId),
					},
				},
			})
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
