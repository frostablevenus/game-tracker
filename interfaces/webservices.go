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

func (handler WebserviceHandler) AddUser(c *gin.Context) (result.Errors, int, result.UserAdd) {
	user := request.User{}
	errors := result.Errors{}
	errMsg := c.BindJSON(&user)
	if errMsg != nil {
		err := result.Error{Message: errMsg}
		errors.Errs = append(errors.Errs, err)
		return errors, 400, result.UserAdd{}
	}

	player := domain.Player{Id: user.PlayerId, Name: user.PlayerName}
	id, errMsg, code := handler.ProfileInteractor.AddUser(player, user.Name, user.Password)
	if errMsg != nil {
		err := result.Error{Message: errMsg}
		errors.Errs = append(errors.Errs, err)
		return errors, code, result.UserAdd{}
	}

	message := result.UserAdd{Id: id, Name: user.Name}
	fmt.Printf("Created user #%d\n", id)
	return errors, 201, message
}

func (handler WebserviceHandler) ShowUser(c *gin.Context) (result.Errors, int, result.User) {
	errors := result.Errors{}
	userId, errMsg := strconv.Atoi(c.Param("id"))
	if errMsg != nil {
		err := result.Error{Message: errMsg, Field: "userId"}
		errors.Errs = append(errors.Errs, err)
		return errors, 400, result.User{}
	}

	name, libraryIds, errMsg, code := handler.ProfileInteractor.ShowUser(userId)
	if errMsg != nil {
		err := result.Error{Message: errMsg}
		errors.Errs = append(errors.Errs, err)
		return errors, code, result.User{}
	}

	var message result.User
	message.Name = name
	message.Id = userId
	for _, libraryId := range libraryIds {
		message.LibraryIds = append(message.LibraryIds, libraryId)
	}
	fmt.Printf("Printed user #%d\n", userId)
	return result.Errors{}, 200, message
}

func (handler WebserviceHandler) RemoveUser(c *gin.Context) (result.Errors, int, result.UserDelete) {
	errors := result.Errors{}
	userId, errMsg := strconv.Atoi(c.Param("id"))
	if errMsg != nil {
		err := result.Error{Message: errMsg, Field: "userId"}
		errors.Errs = append(errors.Errs, err)
		return errors, 400, result.UserDelete{}
	}

	errMsg, code := handler.ProfileInteractor.RemoveUser(userId)
	if errMsg != nil {
		err := result.Error{Message: errMsg}
		errors.Errs = append(errors.Errs, err)
		return errors, code, result.UserDelete{}
	}

	message := result.UserDelete{Id: userId}
	fmt.Printf("Deleted user #%d\n", userId)
	return result.Errors{}, 200, message
}

func (handler WebserviceHandler) ShowUserInfo(c *gin.Context) (result.Errors, int, result.UserInfo) {
	errors := result.Errors{}
	userId, errMsg := strconv.Atoi(c.Param("id"))
	if errMsg != nil {
		err := result.Error{Message: errMsg, Field: "userId"}
		errors.Errs = append(errors.Errs, err)
		return errors, 400, result.UserInfo{}
	}

	info, errMsg, code := handler.ProfileInteractor.ShowUserInfo(userId)
	if errMsg != nil {
		err := result.Error{Message: errMsg}
		errors.Errs = append(errors.Errs, err)
		return errors, code, result.UserInfo{}
	}

	message := result.UserInfo{Id: userId, Info: info}
	fmt.Printf("Printed info of user #%d\n", userId)
	return result.Errors{}, 200, message
}

func (handler WebserviceHandler) EditUserInfo(c *gin.Context) (result.Errors, int, result.UserInfo) {
	errors := result.Errors{}
	userId, errMsg := strconv.Atoi(c.Param("id"))
	if errMsg != nil {
		err := result.Error{Message: errMsg, Field: "userId"}
		errors.Errs = append(errors.Errs, err)
		return errors, 400, result.UserInfo{}
	}
	userInfo := request.UserInfo{}
	errMsg = c.BindJSON(&userInfo)
	if errMsg != nil {
		err := result.Error{Message: errMsg}
		errors.Errs = append(errors.Errs, err)
		return errors, 400, result.UserInfo{}
	}

	errMsg, code := handler.ProfileInteractor.EditUserInfo(userId, userInfo.Info)
	if errMsg != nil {
		err := result.Error{Message: errMsg}
		errors.Errs = append(errors.Errs, err)
		return result.Errors{}, code, result.UserInfo{}
	}

	message := result.UserInfo{Id: userId, Info: userInfo.Info}
	fmt.Printf("Editted info of user #%d\n", userId)
	return result.Errors{}, 200, message
}

func (handler WebserviceHandler) AddLibrary(c *gin.Context) (result.Errors, int, result.LibraryAdd) {
	errors := result.Errors{}
	userId, errMsg := strconv.Atoi(c.Param("id"))
	if errMsg != nil {
		err := result.Error{Message: errMsg, Field: "userId"}
		errors.Errs = append(errors.Errs, err)
		return errors, 400, result.LibraryAdd{}
	}
	id, errMsg, code := handler.ProfileInteractor.AddLibrary(userId)
	if errMsg != nil {
		err := result.Error{Message: errMsg}
		errors.Errs = append(errors.Errs, err)
		return errors, code, result.LibraryAdd{}
	}

	message := result.LibraryAdd{Id: id, UserId: userId}
	fmt.Printf("Added library #%d\n", id)
	return result.Errors{}, 201, message
}

func (handler WebserviceHandler) ShowLibrary(c *gin.Context) (result.Errors, int, result.Library) {
	errors := result.Errors{}
	userId, errMsg := strconv.Atoi(c.Param("id"))
	if errMsg != nil {
		err := result.Error{Message: errMsg, Field: "userId"}
		errors.Errs = append(errors.Errs, err)
		return errors, 400, result.Library{}
	}
	libraryId, errMsg := strconv.Atoi(c.Param("libId"))
	if errMsg != nil {
		err := result.Error{Message: errMsg, Field: "libId"}
		errors.Errs = append(errors.Errs, err)
		return errors, 400, result.Library{}
	}

	gameIds, errMsg, code := handler.ProfileInteractor.ShowLibrary(userId, libraryId)
	if errMsg != nil {
		err := result.Error{Message: errMsg}
		errors.Errs = append(errors.Errs, err)
		return errors, code, result.Library{}
	}

	var message result.Library
	message.Id = libraryId
	message.UserId = userId
	for _, gameId := range gameIds {
		message.GamesIds = append(message.GamesIds, gameId)
	}
	fmt.Printf("Printed library #%d\n", libraryId)
	return result.Errors{}, 200, message
}

func (handler WebserviceHandler) RemoveLibrary(c *gin.Context) (result.Errors, int, result.LibraryDelete) {
	errors := result.Errors{}
	userId, errMsg := strconv.Atoi(c.Param("id"))
	if errMsg != nil {
		err := result.Error{Message: errMsg, Field: "userId"}
		errors.Errs = append(errors.Errs, err)
		return errors, 400, result.LibraryDelete{}
	}
	libraryId, errMsg := strconv.Atoi(c.Param("libId"))
	if errMsg != nil {
		err := result.Error{Message: errMsg, Field: "libId"}
		errors.Errs = append(errors.Errs, err)
		return errors, 400, result.LibraryDelete{}
	}

	errMsg, code := handler.ProfileInteractor.RemoveLibrary(userId, libraryId)
	if errMsg != nil {
		err := result.Error{Message: errMsg}
		errors.Errs = append(errors.Errs, err)
		return errors, code, result.LibraryDelete{}
	}

	message := result.LibraryDelete{Id: libraryId}
	fmt.Printf("Deleted library #%d\n", libraryId)
	return result.Errors{}, 200, message
}

func (handler WebserviceHandler) AddGame(c *gin.Context) (result.Errors, int, result.Game) {
	errors := result.Errors{}
	userId, errMsg := strconv.Atoi(c.Param("id"))
	if errMsg != nil {
		err := result.Error{Message: errMsg, Field: "userId"}
		errors.Errs = append(errors.Errs, err)
		return errors, 400, result.Game{}
	}
	libraryId, errMsg := strconv.Atoi(c.Param("libId"))
	if errMsg != nil {
		err := result.Error{Message: errMsg, Field: "userId"}
		errors.Errs = append(errors.Errs, err)
		return errors, 400, result.Game{}
	}
	game := request.Game{}
	errMsg = c.BindJSON(&game)
	if errMsg != nil {
		err := result.Error{Message: errMsg}
		errors.Errs = append(errors.Errs, err)
		return errors, 400, result.Game{}
	}

	id, errMsg, code := handler.ProfileInteractor.AddGame(userId, libraryId, game.Name, game.Producer, game.Value)
	if errMsg != nil {
		err := result.Error{Message: errMsg}
		errors.Errs = append(errors.Errs, err)
		return errors, code, result.Game{}
	}

	message := result.Game{Id: id, LibraryId: libraryId, UserId: userId, Name: game.Name,
		Producer: game.Producer, Value: game.Value}
	fmt.Printf("Added game #%d\n", id)
	return result.Errors{}, 201, message
}

func (handler WebserviceHandler) PickGame(c *gin.Context) (result.Errors, int, result.GameToLib) {
	errors := result.Errors{}
	userId, errMsg := strconv.Atoi(c.Param("id"))
	if errMsg != nil {
		err := result.Error{Message: errMsg, Field: "userId"}
		errors.Errs = append(errors.Errs, err)
		return errors, 400, result.GameToLib{}
	}
	libraryId, errMsg := strconv.Atoi(c.Param("libId"))
	if errMsg != nil {
		err := result.Error{Message: errMsg, Field: "libId"}
		errors.Errs = append(errors.Errs, err)
		return errors, 400, result.GameToLib{}
	}
	gameId, errMsg := strconv.Atoi(c.Param("gameId"))
	if errMsg != nil {
		err := result.Error{Message: errMsg, Field: "gameId"}
		errors.Errs = append(errors.Errs, err)
		return errors, 400, result.GameToLib{}
	}

	errMsg, code := handler.ProfileInteractor.PickGame(userId, libraryId, gameId)
	if errMsg != nil {
		err := result.Error{Message: errMsg}
		errors.Errs = append(errors.Errs, err)
		return errors, code, result.GameToLib{}
	}

	message := result.GameToLib{Id: gameId, LibraryId: libraryId, UserId: userId}
	fmt.Printf("Added game #%d\n", gameId)
	return result.Errors{}, 201, message
}

func (handler WebserviceHandler) RemoveGame(c *gin.Context) (result.Errors, int, result.GameToLib) {
	errors := result.Errors{}
	userId, errMsg := strconv.Atoi(c.Param("id"))
	if errMsg != nil {
		err := result.Error{Message: errMsg, Field: "userId"}
		errors.Errs = append(errors.Errs, err)
		return errors, 400, result.GameToLib{}
	}
	libraryId, errMsg := strconv.Atoi(c.Param("libId"))
	if errMsg != nil {
		err := result.Error{Message: errMsg, Field: "libId"}
		errors.Errs = append(errors.Errs, err)
		return errors, 400, result.GameToLib{}
	}
	gameId, errMsg := strconv.Atoi(c.Param("gameId"))
	if errMsg != nil {
		err := result.Error{Message: errMsg, Field: "gameId"}
		errors.Errs = append(errors.Errs, err)
		return errors, 400, result.GameToLib{}
	}

	errMsg, code := handler.ProfileInteractor.RemoveGame(userId, libraryId, gameId)
	if errMsg != nil {
		err := result.Error{Message: errMsg}
		errors.Errs = append(errors.Errs, err)
		return errors, code, result.GameToLib{}
	}

	message := result.GameToLib{Id: gameId, LibraryId: libraryId}
	fmt.Printf("Deleted game #%d\n", gameId)
	return result.Errors{}, 200, message
}

func (handler WebserviceHandler) GetViewError(errs result.Errors) []gin.H {
	errors := []gin.H{}
	for _, err := range errs.Errs {
		if err.Field != "" {
			errors = append(errors, gin.H{
				"field":   err.Field,
				"message": err.Message.Error(),
			})
		} else {
			errors = append(errors, gin.H{
				"message": err.Message.Error(),
			})
		}
	}
	return errors
}
