package interfaces

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"game-tracker/domain"
	"game-tracker/usecases"
)

type ProfileInteractor interface {
	ShowLibrary(userId, libraryId int) (string, []domain.Game, error)
	AddUser(player domain.Player, userName string) error
	RemoveUser(playerId, userId int) error
	EditUserInfo(userId int, info string) error
	AddGame(userId, libraryId int, gameName, gameProducer string, gameValue float64) error
	RemoveGame(userId, libraryId, gameId int) error
}

type WebserviceHandler struct {
	ProfileInteractor usecases.ProfileInteractor
}

func (handler WebserviceHandler) ShowLibrary(res http.ResponseWriter, req *http.Request) error {
	userId, err := getFormUserId(req)
	if err != nil {
		return err
	}
	libraryId, err := getFormLibraryId(req)
	if err != nil {
		return err
	}

	info, games, err := handler.ProfileInteractor.ShowLibrary(userId, libraryId)
	if err != nil {
		return err
	}

	io.WriteString(res, fmt.Sprintf("Library #%d of user #%d\n", libraryId, userId))
	io.WriteString(res, fmt.Sprintf("User information: %s\n", info))
	io.WriteString(res, fmt.Sprintf("Games: \n"))
	for _, game := range games {
		io.WriteString(res, fmt.Sprintf("game id: %d\n", game.Id))
		io.WriteString(res, fmt.Sprintf("game name: %v\n", game.Name))
		io.WriteString(res, fmt.Sprintf("game producer: %v\n", game.Producer))
		io.WriteString(res, fmt.Sprintf("game value: %s\n\n", game.Value))
	}
	fmt.Printf("Printed library #%d of user #%d\n", libraryId, userId)
	return nil
}

func (handler WebserviceHandler) AddUser(res http.ResponseWriter, req *http.Request) error {
	playerId, err := getFormPlayerId(req)
	if err != nil {
		return err
	}
	playerName, err := getFormPlayerName(req)
	if err != nil {
		return err
	}
	userName, err := getFormUserName(req)
	if err != nil {
		return err
	}

	player := domain.Player{Id: playerId, Name: playerName}
	err = handler.ProfileInteractor.AddUser(player, userName)
	if err != nil {
		return err
	}

	io.WriteString(res, fmt.Sprintf(
		"Player '%s' (id #%d) created account with username: %s\n",
		player.Name, player.Id, userName))
	return nil
}

func (handler WebserviceHandler) RemoveUser(res http.ResponseWriter, req *http.Request) error {
	userId, err := getFormUserId(req)
	if err != nil {
		return err
	}
	playerId, err := getFormPlayerId(req)
	if err != nil {
		return err
	}

	err = handler.ProfileInteractor.RemoveUser(playerId, userId)
	if err != nil {
		return err
	}

	io.WriteString(res, fmt.Sprintf("Player #%d deleted user account #%d\n", playerId, userId))
	return nil
}

func (handler WebserviceHandler) EditUserInfo(res http.ResponseWriter, req *http.Request) error {
	userId, err := getFormUserId(req)
	if err != nil {
		return err
	}
	targetId, err := getFormTargetId(req)
	if err != nil {
		return err
	}
	info := req.FormValue("info")
	err = handler.ProfileInteractor.EditUserInfo(userId, targetId, info)
	if err != nil {
		return err
	}
	io.WriteString(res, fmt.Sprintf("Added personal information for user #%d\n", userId))
	return nil
}

func (handler WebserviceHandler) AddLibrary(res http.ResponseWriter, req *http.Request) error {
	userId, err := getFormUserId(req)
	if err != nil {
		return err
	}
	err = handler.ProfileInteractor.AddLibrary(userId)
	if err != nil {
		return err
	}
	io.WriteString(res, fmt.Sprintf(
		"User #%d added library", userId))
	return nil
}

func (handler WebserviceHandler) RemoveLibrary(res http.ResponseWriter, req *http.Request) error {
	userId, err := getFormUserId(req)
	if err != nil {
		return err
	}
	libraryId, err := getFormLibraryId(req)
	if err != nil {
		return err
	}

	err = handler.ProfileInteractor.RemoveLibrary(userId, libraryId)
	if err != nil {
		return err
	}

	io.WriteString(res, fmt.Sprintf(
		"User #%d removed library #%d", userId, libraryId))
	return nil
}

func (handler WebserviceHandler) AddGame(res http.ResponseWriter, req *http.Request) error {
	userId, err := getFormUserId(req)
	if err != nil {
		return err
	}
	libraryId, err := getFormLibraryId(req)
	if err != nil {
		return err
	}
	gameName, err := getFormGameName(req)
	if err != nil {
		return err
	}
	gameProducer, err := getFormGameProducer(req)
	if err != nil {
		return err
	}
	gameValue, err := getFormGameValue(req)
	if err != nil {
		return err
	}

	err = handler.ProfileInteractor.AddGame(userId, libraryId, gameName, gameProducer, gameValue)
	if err != nil {
		return err
	}

	io.WriteString(res, fmt.Sprintf(
		"User #%d added game to library #%d:\nGame name: %s\nGame producer: %s\nGame value: %s\n",
		userId, libraryId, gameName, gameProducer, gameValue))
	return nil
}

func (handler WebserviceHandler) RemoveGame(res http.ResponseWriter, req *http.Request) error {
	userId, err := getFormUserId(req)
	if err != nil {
		return err
	}
	libraryId, err := getFormLibraryId(req)
	if err != nil {
		return err
	}
	gameId, err := getFormUserId(req)
	if err != nil {
		return err
	}

	err = handler.ProfileInteractor.RemoveGame(userId, libraryId, gameId)
	if err != nil {
		return err
	}
	io.WriteString(res, fmt.Sprintf(
		"User #%d removed game (id #%d) from library #%d\n",
		userId, gameId, libraryId))
	return nil
}

func getFormPlayerId(req *http.Request) (int, error) {
	var form string
	if form = req.FormValue("playerId"); form == "" {
		err := fmt.Errorf("playerId cannot be empty")
		return 0, err
	}
	playerId, err := strconv.Atoi(form)
	return playerId, err
}

func getFormUserId(req *http.Request) (int, error) {
	var form string
	if form = req.FormValue("userId"); form == "" {
		err := fmt.Errorf("userId cannot be empty")
		return 0, err
	}
	userId, err := strconv.Atoi(form)
	return userId, err
}

func getFormTargetId(req *http.Request) (int, error) {
	var form string
	if form = req.FormValue("targetId"); form == "" {
		err := fmt.Errorf("targetId cannot be empty")
		return 0, err
	}
	targetId, err := strconv.Atoi(form)
	return targetId, err
}

func getFormLibraryId(req *http.Request) (int, error) {
	var form string
	if form = req.FormValue("libraryId"); form == "" {
		err := fmt.Errorf("libraryId cannot be empty")
		return 0, err
	}
	libraryId, err := strconv.Atoi(form)
	return libraryId, err
}

func getFormGameId(req *http.Request) (int, error) {
	var form string
	if form = req.FormValue("gameId"); form == "" {
		err := fmt.Errorf("gameId cannot be empty")
		return 0, err
	}
	gameId, err := strconv.Atoi(form)
	return gameId, err
}

func getFormPlayerName(req *http.Request) (string, error) {
	playerName := req.FormValue("playerName")
	if playerName == "" {
		err := fmt.Errorf("playerName cannot be empty")
		return "", err
	}
	return playerName, nil
}

func getFormUserName(req *http.Request) (string, error) {
	userName := req.FormValue("userName")
	if userName == "" {
		err := fmt.Errorf("userName cannot be empty")
		return "", err
	}
	return userName, nil
}

func getFormGameName(req *http.Request) (string, error) {
	gameName := req.FormValue("gameName")
	if gameName == "" {
		err := fmt.Errorf("gameName cannot be empty")
		return "", err
	}
	return gameName, nil
}

func getFormGameProducer(req *http.Request) (string, error) {
	gameProducer := req.FormValue("gameProducer")
	if gameProducer == "" {
		err := fmt.Errorf("gameProducer cannot be empty")
		return "", err
	}
	return gameProducer, nil
}

func getFormGameValue(req *http.Request) ([]uint8, error) {
	var form string
	if form = req.FormValue("gameValue"); form == "" {
		err := fmt.Errorf("gameValue cannot be empty")
		return []uint8{}, err
	}
	gameValue := []uint8(form)
	return gameValue, nil
}
