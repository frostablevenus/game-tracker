package interfaces

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"strconv"

	"game-tracker/domain"
	"game-tracker/models/request"
	"game-tracker/usecases"
)

type ProfileInteractor interface {
	AddUser(player domain.Player, userName string) (int, error, int)
	ShowUser(userId int) ([]int, error, int)
	RemoveUser(playerId, userId int) (error, int)
	ShowUserInfo(userId int) (string, error, int)
	EditUserInfo(userId int, info string) (error, int)
	AddLibrary(userId int) (error, int)
	ShowLibrary(userId, libraryId int) ([]domain.Game, error, int)
	RemoveLibrary(userId, libraryId int) (error, int)
	AddGame(userId, libraryId int, gameName, gameProducer string, gameValue float64) (int, error, int)
	RemoveGame(userId, libraryId, gameId int) (error, int)
}

type WebserviceHandler struct {
	ProfileInteractor usecases.ProfileInteractor
}

func (handler WebserviceHandler) AddUser(c *gin.Context) (error, int) {
	user := request.User{}
	err := c.BindJSON(&user)
	if err != nil {
		return err, 400
	}

	player := domain.Player{Id: user.PlayerId, Name: user.PlayerName}
	id, err, code := handler.ProfileInteractor.AddUser(player, user.Name)
	if err != nil {
		return err, code
	}

	_, err = io.WriteString(c.Writer, fmt.Sprintf("Player '%s' (id #%d) created user #%d with username: '%s\n'",
		player.Name, player.Id, id, user.Name))
	if err != nil {
		return err, 500
	}
	return nil, 201
}

func (handler WebserviceHandler) ShowUser(c *gin.Context) (error, int) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err, 400
	}

	libraryIds, err, code := handler.ProfileInteractor.ShowUser(userId)
	if err != nil {
		return err, code
	}
	_, err = io.WriteString(c.Writer, fmt.Sprintf("User #%d\nLibraries IDs:\n", userId))
	if err != nil {
		return err, 500
	}
	for _, libraryId := range libraryIds {
		_, err = io.WriteString(c.Writer, fmt.Sprintf("#%d ", libraryId))
		if err != nil {
			return err, 500
		}
	}
	fmt.Printf("Printed user #%d\n", userId)
	return nil, 200
}

func (handler WebserviceHandler) RemoveUser(c *gin.Context) (error, int) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err, 400
	}

	reqId, exist := c.GetQuery("playerId")
	if !exist {
		err = fmt.Errorf("playerId must not be empty")
		return err, 400
	}
	playerId, err := strconv.Atoi(reqId)
	if err != nil {
		return err, 400
	}

	err, code := handler.ProfileInteractor.RemoveUser(playerId, userId)
	if err != nil {
		return err, code
	}

	_, err = io.WriteString(c.Writer, fmt.Sprintf("Player #%d deleted user account #%d\n", playerId, userId))
	if err != nil {
		return err, 500
	}
	return nil, 200
}

func (handler WebserviceHandler) ShowUserInfo(c *gin.Context) (error, int) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err, 400
	}

	info, err, code := handler.ProfileInteractor.ShowUserInfo(userId)
	if err != nil {
		return err, code
	}

	_, err = io.WriteString(c.Writer, fmt.Sprintf("Information of user #%d:\n%s", userId, info))
	if err != nil {
		return err, 500
	}
	return nil, 200
}

func (handler WebserviceHandler) EditUserInfo(c *gin.Context) (error, int) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err, 400
	}
	info := request.UserInfo{}
	err = c.BindJSON(&info)
	if err != nil {
		return err, 400
	}

	err, code := handler.ProfileInteractor.EditUserInfo(userId, info.Info)
	if err != nil {
		return err, code
	}
	_, err = io.WriteString(c.Writer, fmt.Sprintf("Added personal information for user #%d\n", userId))
	if err != nil {
		return err, 500
	}
	return nil, 200
}

func (handler WebserviceHandler) AddLibrary(c *gin.Context) (error, int) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err, 400
	}
	id, err, code := handler.ProfileInteractor.AddLibrary(userId)
	if err != nil {
		return err, code
	}
	_, err = io.WriteString(c.Writer, fmt.Sprintf("User #%d added library #%d", userId, id))
	if err != nil {
		return err, 500
	}
	return nil, 201
}

func (handler WebserviceHandler) ShowLibrary(c *gin.Context) (error, int) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err, 400
	}
	libraryId, err := strconv.Atoi(c.Param("libId"))
	if err != nil {
		return err, 400
	}

	games, err, code := handler.ProfileInteractor.ShowLibrary(userId, libraryId)
	if err != nil {
		return err, code
	}
	_, err = io.WriteString(c.Writer, fmt.Sprintf("Library #%d of user #%d:\n",
		libraryId, userId))
	if err != nil {
		return err, 500
	}
	for _, game := range games {
		message := "Id: %d\nName: %v\nProducer: %v\nValue: %.2f\n\n"
		_, err = io.WriteString(c.Writer, fmt.Sprintf(message, game.Id, game.Name, game.Producer, game.Value))
		if err != nil {
			return err, 500
		}
	}

	fmt.Printf("Printed library #%d of user #%d\n", libraryId, userId)
	return nil, 200
}

func (handler WebserviceHandler) RemoveLibrary(c *gin.Context) (error, int) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err, 400
	}
	libraryId, err := strconv.Atoi(c.Param("libId"))
	if err != nil {
		return err, 400
	}

	err, code := handler.ProfileInteractor.RemoveLibrary(userId, libraryId)
	if err != nil {
		return err, code
	}

	_, err = io.WriteString(c.Writer, fmt.Sprintf("User #%d removed library #%d", userId, libraryId))
	if err != nil {
		return err, 500
	}
	return nil, 200
}

func (handler WebserviceHandler) AddGame(c *gin.Context) (error, int) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err, 400
	}
	libraryId, err := strconv.Atoi(c.Param("libId"))
	if err != nil {
		return err, 400
	}
	game := request.Game{}
	err = c.BindJSON(&game)
	if err != nil {
		return err, 400
	}

	id, err, code := handler.ProfileInteractor.AddGame(userId, libraryId, game.Name, game.Producer, game.Value)
	if err != nil {
		return err, code
	}

	_, err = io.WriteString(c.Writer, fmt.Sprintf(
		"User #%d added game #%d to library #%d:\nGame name: %s\nGame producer: %s\nGame value: %.2f\n",
		userId, id, libraryId, game.Name, game.Producer, game.Value))
	if err != nil {
		return err, 500
	}
	return nil, 201
}

func (handler WebserviceHandler) RemoveGame(c *gin.Context) (error, int) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err, 400
	}
	libraryId, err := strconv.Atoi(c.Param("libId"))
	if err != nil {
		return err, 400
	}
	gameId, err := strconv.Atoi(c.Param("gameId"))
	if err != nil {
		return err, 400
	}

	err, code := handler.ProfileInteractor.RemoveGame(userId, libraryId, gameId)
	if err != nil {
		return err, code
	}
	_, err = io.WriteString(c.Writer, fmt.Sprintf("User #%d removed game (id #%d) from library #%d\n",
		userId, gameId, libraryId))
	if err != nil {
		return err, 500
	}
	return nil, 200
}
