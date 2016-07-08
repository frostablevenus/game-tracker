package usecases

import (
	"fmt"

	"game-tracker/domain"
)

type UserRepository interface {
	Store(user User) error
	Remove(user User) error
	FindById(id int) (User, error)
	Count() int
	IsUnique(userName string) bool
	StoreInfo(user User, info string)
	LoadInfo(user User) string
}

type LibraryRepository interface {
	Store(library Library) error
	FindById(id int) (Library, error)
}

type User struct {
	Id           int
	Name         string
	Player       domain.Player //This user (account) was created by some player
	PersonalInfo string
}

type Library struct {
	Id     int
	Player domain.Player //This library belongs to some player
	Games  []Game
}

type Game struct {
	Id       int
	Name     string
	Producer string
	Value    float64
}

type Logger interface {
	Log(message string) error
}

type ProfileInteractor struct {
	UserRepository    UserRepository
	LibraryRepository LibraryRepository
	GameRepository    domain.GameRepository
	Logger            Logger
}

func (interactor *ProfileInteractor) ShowLibrary(userId, libraryId int) (string, []Game, error) {
	var games []Game
	user, err := interactor.UserRepository.FindById(userId)
	if err != nil {
		fmt.Println("User #%i does not exist", userId)
		return "", nil, err
	}

	info := interactor.UserRepository.LoadInfo(user)

	library, err := interactor.LibraryRepository.FindById(libraryId)
	if err != nil {
		fmt.Println("Library #%i of user #%i does not exist", libraryId, userId)
		return "", nil, err
	}
	if user.Player.Id != library.Player.Id {
		message := "User #%i (player #%i) "
		message += "is not allowed to see games "
		message += "in library #%i (of another player #%i)"
		err := fmt.Errorf(message,
			user.Id,
			user.Player.Id,
			library.Id,
			library.Player.Id)
		interactor.Logger.Log(err.Error())
		games = make([]Game, 0)
		return info, games, err
	} else {
		games = make([]Game, len(library.Games))
		for i, game := range library.Games {
			games[i] = Game{game.Id, game.Name, game.Producer, game.Value}
		}
		return info, games, nil
	}

}

func (interactor *ProfileInteractor) AddUser(player domain.Player, userName string) error {
	user := User{interactor.UserRepository.Count(), userName, player, ""}

	// Application rule: usernames cannot repeat
	if !interactor.UserRepository.IsUnique(userName) {
		err := fmt.Errorf("Username #%s is taken", userName)
		interactor.Logger.Log(err.Error())
		return err
	}

	err := interactor.UserRepository.Store(user)
	if err != nil {
		interactor.Logger.Log(err.Error())
		return err
	}
	interactor.Logger.Log("Added user #%i (player #%i)")
	return nil
}

func (interactor *ProfileInteractor) RemoveUser(playerId, userId int) error {
	player, err := interactor.UserRepository.FindById(playerId)
	if err != nil {
		interactor.Logger.Log(err.Error())
		return err
	}

	user, err := interactor.UserRepository.FindById(userId)
	if err != nil {
		interactor.Logger.Log(err.Error())
		return err
	}

	if player.Id != user.Player.Id {
		err := fmt.Errorf("Player #%i cannot remove user account of player #%i",
			playerId, user.Player.Id)
		interactor.Logger.Log(err.Error())
		return err
	}

	interactor.UserRepository.Remove(user)
	interactor.Logger.Log(fmt.Sprintf("Removed user #%s (id #%i)", user.Name, user.Id))
	return nil
}

func (interactor *ProfileInteractor) EditUserInfo(userId int, info string) error {
	user, err := interactor.UserRepository.FindById(userId)
	if err != nil {
		interactor.Logger.Log(err.Error())
		return err
	}
	interactor.UserRepository.StoreInfo(user, info)
	interactor.Logger.Log(fmt.Sprintf("Editted information of user #%s (id #%i)", user.Name, user.Id))
	return nil
}

func (interactor *ProfileInteractor) AddGame(userId, libraryId int, gameName, gameProducer string, gameValue float64) error {
	user, err := interactor.UserRepository.FindById(userId)
	if err != nil {
		interactor.Logger.Log(err.Error())
		return err
	}
	library, err := interactor.LibraryRepository.FindById(libraryId)
	if err != nil {
		interactor.Logger.Log(err.Error())
		return err
	}
	message := ""
	if user.Player.Id != library.Player.Id {
		message = "User #%i (player #%i) "
		message += "is not allowed to add games "
		message += "to library #%i (of another player #%i)"
		err := fmt.Errorf(message,
			user.Id,
			user.Player.Id,
			library.Id,
			library.Player.Id)
		interactor.Logger.Log(err.Error())
		return err
	}

	gameId := len(library.Games)
	game := Game{gameId, gameName, gameProducer, gameValue}
	library.Games = append(library.Games, game)
	err = interactor.LibraryRepository.Store(library)
	if err != nil {
		interactor.Logger.Log(err.Error())
		return err
	}

	interactor.Logger.Log(fmt.Sprintf(
		"User added game '%s' (id #%i) to library #%i",
		game.Name, game.Id, library.Id))
	return nil
}

func (interactor *ProfileInteractor) RemoveGame(userId, libraryId, gameId int) error {
	user, err := interactor.UserRepository.FindById(userId)
	if err != nil {
		interactor.Logger.Log(err.Error())
		return err
	}
	library, err := interactor.LibraryRepository.FindById(libraryId)
	if err != nil {
		interactor.Logger.Log(err.Error())
		return err
	}
	message := ""
	if user.Player.Id != library.Player.Id {
		message = "User #%i (player #%i) "
		message += "is not allowed to remove games "
		message += "from library #%i (of another player #%i)"
		err := fmt.Errorf(message,
			user.Id,
			user.Player.Id,
			library.Id,
			library.Player.Id)
		interactor.Logger.Log(err.Error())
		return err
	}

	game, err := interactor.GameRepository.FindById(gameId)
	if err != nil {
		interactor.Logger.Log(err.Error())
		return err
	}

	for i := range library.Games {
		if game.Id == library.Games[i].Id {
			library.Games = append(library.Games[:i], library.Games[i+1:]...)
			err = interactor.LibraryRepository.Store(library)
			if err != nil {
				return err
			}
			interactor.Logger.Log(fmt.Sprintf(
				"User removed game '%s' (id #%i) from library #%i",
				game.Name, game.Id, library.Id))
			return nil
		}
	}
	message = "Library #%i of user #%i (player #%i) does not contain game #%i"
	err = fmt.Errorf(message, library.Id, user.Id, user.Player.Id, game.Id)
	return err
}
