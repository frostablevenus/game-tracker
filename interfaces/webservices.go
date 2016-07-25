package interfaces

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"

	"game-tracker/domain"
	"game-tracker/models/request"
	"game-tracker/models/result"
	"game-tracker/usecases"
)

type ProfileInteractor interface {
	AddUser(player domain.Player, userName, password string) (int, error, int)
	ShowUser(userId int) (string, []int, error, int)
	RemoveUser(userId int) (error, int)
	ShowUserInfo(userId int) (string, error, int)
	EditUserInfo(userId int, info string) (error, int)
	AddLibrary(userId int) (error, int)
	ShowLibrary(userId, libraryId int) ([]domain.Game, error, int)
	RemoveLibrary(userId, libraryId int) (error, int)
	AddGame(userId, libraryId int, gameName, gameProducer string, gameValue float64) (int, error, int)
	RemoveGame(userId, libraryId, gameId int) (error, int)
	FindLoginId(username, password string) (int, error, int)
}

type WebserviceHandler struct {
	ProfileInteractor usecases.ProfileInteractor
}

func (handler WebserviceHandler) AddUser(c *gin.Context) (error, int, string) {
	user := request.User{}
	err := c.BindJSON(&user)
	if err != nil {
		return err, 400, ""
	}

	player := domain.Player{Id: user.PlayerId, Name: user.PlayerName}
	id, err, code := handler.ProfileInteractor.AddUser(player, user.Name, user.Password)
	if err != nil {
		return err, code, ""
	}

	message := fmt.Sprintf("Player '%s' (id #%d) created user #%d with username: '%s'\n",
		player.Name, player.Id, id, user.Name)
	fmt.Printf("Created user #%d\n", id)
	return nil, 201, message
}

func (handler WebserviceHandler) ShowUser(c *gin.Context) (error, int, result.User) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err, 400, result.User{}
	}

	name, libraryIds, err, code := handler.ProfileInteractor.ShowUser(userId)
	if err != nil {
		return err, code, result.User{}
	}

	var message result.User
	message.Name = name
	message.Id = userId
	for _, libraryId := range libraryIds {
		message.LibraryIds = append(message.LibraryIds, libraryId)
	}
	fmt.Printf("Printed user #%d\n", userId)
	return nil, 200, message
}

func (handler WebserviceHandler) RemoveUser(c *gin.Context) (error, int, string) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err, 400, ""
	}

	err, code := handler.ProfileInteractor.RemoveUser(userId)
	if err != nil {
		return err, code, ""
	}

	message := fmt.Sprintf("Removed user account #%d\n", userId)
	fmt.Printf("Deleted user #%d\n", userId)
	return nil, 200, message
}

func (handler WebserviceHandler) ShowUserInfo(c *gin.Context) (error, int, result.UserInfo) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err, 400, result.UserInfo{}
	}

	info, err, code := handler.ProfileInteractor.ShowUserInfo(userId)
	if err != nil {
		return err, code, result.UserInfo{}
	}

	var message result.UserInfo
	message.Id = userId
	message.Info = info
	fmt.Printf("Printed info of user #%d\n", userId)
	return nil, 200, message
}

func (handler WebserviceHandler) EditUserInfo(c *gin.Context) (error, int, string) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err, 400, ""
	}
	info := request.UserInfo{}
	err = c.BindJSON(&info)
	if err != nil {
		return err, 400, ""
	}

	err, code := handler.ProfileInteractor.EditUserInfo(userId, info.Info)
	if err != nil {
		return err, code, ""
	}

	message := fmt.Sprintf("Editted personal information for user #%d\n", userId)
	fmt.Printf("Editted info of user #%d\n", userId)
	return nil, 200, message
}

func (handler WebserviceHandler) AddLibrary(c *gin.Context) (error, int, string) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err, 400, ""
	}
	id, err, code := handler.ProfileInteractor.AddLibrary(userId)
	if err != nil {
		return err, code, ""
	}

	message := fmt.Sprintf("User #%d added library #%d", userId, id)
	fmt.Printf("Added library #%d\n", id)
	return nil, 201, message
}

func (handler WebserviceHandler) ShowLibrary(c *gin.Context) (error, int, result.Library) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err, 400, result.Library{}
	}
	libraryId, err := strconv.Atoi(c.Param("libId"))
	if err != nil {
		return err, 400, result.Library{}
	}

	games, err, code := handler.ProfileInteractor.ShowLibrary(userId, libraryId)
	if err != nil {
		return err, code, result.Library{}
	}

	var message result.Library
	message.Id = libraryId
	message.UserId = userId
	for _, game := range games {
		message.Games = append(message.Games, result.Game{
			Id:       game.Id,
			Name:     game.Name,
			Producer: game.Producer,
			Value:    game.Value})
	}
	fmt.Printf("Printed library #%d\n", libraryId)
	return nil, 200, message
}

func (handler WebserviceHandler) RemoveLibrary(c *gin.Context) (error, int, string) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err, 400, ""
	}
	libraryId, err := strconv.Atoi(c.Param("libId"))
	if err != nil {
		return err, 400, ""
	}

	err, code := handler.ProfileInteractor.RemoveLibrary(userId, libraryId)
	if err != nil {
		return err, code, ""
	}

	message := fmt.Sprintf("User #%d removed library #%d", userId, libraryId)
	fmt.Printf("Deleted library #%d\n", libraryId)
	return nil, 200, message
}

func (handler WebserviceHandler) AddGame(c *gin.Context) (error, int, string) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err, 400, ""
	}
	libraryId, err := strconv.Atoi(c.Param("libId"))
	if err != nil {
		return err, 400, ""
	}
	game := request.Game{}
	err = c.BindJSON(&game)
	if err != nil {
		return err, 400, ""
	}

	id, err, code := handler.ProfileInteractor.AddGame(userId, libraryId, game.Name, game.Producer, game.Value)
	if err != nil {
		return err, code, ""
	}

	message := fmt.Sprintf("User #%d added game #%d to library #%d\n", userId, id, libraryId)
	fmt.Printf("Added game #%d\n", id)
	return nil, 201, message
}

func (handler WebserviceHandler) RemoveGame(c *gin.Context) (error, int, string) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err, 400, ""
	}
	libraryId, err := strconv.Atoi(c.Param("libId"))
	if err != nil {
		return err, 400, ""
	}
	gameId, err := strconv.Atoi(c.Param("gameId"))
	if err != nil {
		return err, 400, ""
	}

	err, code := handler.ProfileInteractor.RemoveGame(userId, libraryId, gameId)
	if err != nil {
		return err, code, ""
	}
	message := fmt.Sprintf("User #%d removed game (id #%d) from library #%d\n",
		userId, gameId, libraryId)
	fmt.Printf("Deleted game #%d\n", gameId)
	return nil, 200, message
}
