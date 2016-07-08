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
	ShowLibrary(userId, libraryId int) (string, []usecases.Game, error)
	AddUser(player domain.Player, userName string) error
	RemoveUser(playerId, userId int) error
	EditUserInfo(userId int, info string) error
	AddGame(userId, libraryId int, gameName, gameProducer string, gameValue float64) error
	RemoveGame(userId, libraryId, gameId int) error
}

type WebserviceHandler struct {
	ProfileInteractor ProfileInteractor //In production stage change this to usecases.ProfileInteractor
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

	info, games, _ := handler.ProfileInteractor.ShowLibrary(userId, libraryId)
	io.WriteString(res, fmt.Sprintf("User information: %s\n", info))

	for _, game := range games {
		io.WriteString(res, fmt.Sprintf("game id: %d\n", game.Id))
		io.WriteString(res, fmt.Sprintf("game name: %v\n", game.Name))
		io.WriteString(res, fmt.Sprintf("game producer: %v\n", game.Producer))
		io.WriteString(res, fmt.Sprintf("game value: %f\n\n", game.Value))
	}

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

	err = handler.ProfileInteractor.AddUser(domain.Player{playerId, playerName}, userName)
	if err != nil {
		return err
	}

	io.WriteString(res, fmt.Sprintf(
		"Player '%s' (id #%i) created account with username: %s\n",
		playerName, playerId, userName))
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

	io.WriteString(res, fmt.Sprintf("Player #%i deleted user account #%i\n", playerId, userId))
	return nil
}

func (handler WebserviceHandler) EditUserInfo(res http.ResponseWriter, req *http.Request) error {
	userId, err := getFormUserId(req)
	if err != nil {
		return err
	}
	info := req.FormValue("info")

	err = handler.ProfileInteractor.EditUserInfo(userId, info)
	if err != nil {
		return err
	}

	io.WriteString(res, fmt.Sprintf("Added personal information for user #%i\n", userId))
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
		"User #%i added game to library #%i:\ngame id: %i\ngame name: %s\ngame producer: %s\ngame value: %v\n",
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
		"User #%i removed game (id #%i) from library #%i\n",
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

func getFormGameValue(req *http.Request) (float64, error) {
	var form string
	if form = req.FormValue("gameValue"); form == "" {
		err := fmt.Errorf("gameValue cannot be empty")
		return 0, err
	}
	gameValue, err := strconv.ParseFloat(form, 64)
	return gameValue, err
}
