package routes

import (
	"fmt"
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
		err, code, message := webserviceHandler.ShowUser(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
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
		err, code, message := webserviceHandler.AddUser(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
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
		err, code, message := webserviceHandler.ShowUserInfo(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
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
		err, code, _ := webserviceHandler.RemoveUser(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.Status(204)
		}
	})
	users.PUT("/info", func(c *gin.Context) {
		err, code, message := webserviceHandler.EditUserInfo(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
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
		err, code, message := webserviceHandler.ShowLibrary(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
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
		err, code, message := webserviceHandler.AddLibrary(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
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
		err, code, _ := webserviceHandler.RemoveLibrary(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.Status(204)
		}
	})

	games := libraries.Group("/:libId/games")
	games.POST("", func(c *gin.Context) {
		err, code, message := webserviceHandler.AddGame(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
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
		err, code, message := webserviceHandler.PickGame(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
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
		err, code, _ := webserviceHandler.RemoveGame(c)
		if err != nil {
			c.AbortWithError(code, err)
		} else {
			c.Status(204)
		}
	})
	return engine
}
