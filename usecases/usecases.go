package usecases

import (
	"fmt"

	"game-tracker/src/domain"
)

type UserRepository interface {
	Store(user User)
	FindById(id int) (User, error)
	StoreInfo(user User, info string)
	LoadInfo(user User) string
}

type LibraryRepository interface {
	Store(library Library)
	FindById(id int) (Library, error)
}

type User struct {
	Id           int
	PersonalInfo string
	Player       domain.Player
}

type Library struct {
	Id     int
	Player domain.Player //belongs to player A
	Games  []domain.Game
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

func (interactor *ProfileInteractor) Profile(userId, libraryId int) (string, []domain.Game, error) {
	var games []domain.Game
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
		games = make([]domain.Game, 0)
		return info, games, err
	} else {
		games = make([]domain.Game, len(library.Games))
		for i, game := range library.Games {
			games[i] = domain.Game{game.Id, game.Name, game.Producer, game.Value}
		}
		return info, games, nil
	}

}

func (interactor *ProfileInteractor) EditUserInfo(userId int, info string) error {
	user, err := interactor.UserRepository.FindById(userId)
	if err != nil {
		fmt.Println("User #%i does not exist", userId)
		return err
	}
	interactor.UserRepository.StoreInfo(user, info)
	return nil
}

func (interactor *ProfileInteractor) Add(userId, libraryId, gameId int) error {
	user, err := interactor.UserRepository.FindById(userId)
	if err != nil {
		fmt.Println("User #%i does not exist", userId)
		return err
	}
	library, err := interactor.LibraryRepository.FindById(libraryId)
	if err != nil {
		fmt.Println("Library #%i of user #%i does not exist", libraryId, userId)
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

	game, err := interactor.GameRepository.FindById(gameId)
	library.Games = append(library.Games, game)
	interactor.LibraryRepository.Store(library)
	interactor.Logger.Log(fmt.Sprintf(
		"User added game '%s' (#%i) to library #%i",
		game.Name, game.Id, library.Id))
	return nil
}

func (interactor *ProfileInteractor) Remove(userId, libraryId, gameId int) error {
	user, err := interactor.UserRepository.FindById(userId)
	library, err := interactor.LibraryRepository.FindById(libraryId)
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
	for i := range library.Games {
		if game.Id == library.Games[i].Id {
			library.Games = append(library.Games[:i], library.Games[i+1:]...)
			interactor.LibraryRepository.Store(library)
			return nil
		}
	}
	message = "Library #%i of user #%i (player #%i) does not contain game #%i"
	err = fmt.Errorf(message, library.Id, user.Id, user.Player.Id, game.Id)
	return err
}
