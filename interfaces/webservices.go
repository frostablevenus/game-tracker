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
	AddLibrary(userId int) (int, error)
	RemoveLibrary(userId, libraryId int) error
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

	message := "Library #%d of user #%d\nUser information: %s\nGames: \n"
	_, err = io.WriteString(res, fmt.Sprintf(message, libraryId, userId, info))
	if err != nil {
		return err
	}
	for _, game := range games {
		message = "game id: %d\ngame name: %v\ngame producer: %v\ngame value: %s\n\n"
		_, err = io.WriteString(res, fmt.Sprintf(message, game.Id, game.Name, game.Producer, game.Value))
		if err != nil {
			return err
		}
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
	id, err := handler.ProfileInteractor.AddUser(player, userName)
	if err != nil {
		return err
	}

	_, err = io.WriteString(res, fmt.Sprintf("Player '%s' (id #%d) created user #%d with username: '%s\n'",
		player.Name, player.Id, id, userName))
	return err
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

	_, err = io.WriteString(res, fmt.Sprintf("Player #%d deleted user account #%d\n", playerId, userId))
	return err
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
	_, err = io.WriteString(res, fmt.Sprintf("Added personal information for user #%d\n", userId))
	return err
}

func (handler WebserviceHandler) AddLibrary(res http.ResponseWriter, req *http.Request) error {
	userId, err := getFormUserId(req)
	if err != nil {
		return err
	}
	id, err := handler.ProfileInteractor.AddLibrary(userId)
	if err != nil {
		return err
	}
	_, err = io.WriteString(res, fmt.Sprintf("User #%d added library #%d", userId, id))
	return err
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

	_, err = io.WriteString(res, fmt.Sprintf("User #%d removed library #%d", userId, libraryId))
	return err
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

	_, err = handler.ProfileInteractor.AddGame(userId, libraryId, gameName, gameProducer, gameValue)
	if err != nil {
		return err
	}

	_, err = io.WriteString(res, fmt.Sprintf(
		"User #%d added game to library #%d:\nGame name: %s\nGame producer: %s\nGame value: %s\n",
		userId, libraryId, gameName, gameProducer, gameValue))
	return err
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
	_, err = io.WriteString(res, fmt.Sprintf("User #%d removed game (id #%d) from library #%d\n",
		userId, gameId, libraryId))
	return err
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
