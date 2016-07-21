package usecases

import (
	"fmt"

	"game-tracker/domain"
)

type UserRepository interface {
	Store(user User) (int, error)
	Remove(user User) error
	FindById(id int) (User, error, int)
	UserExisted(userName string) (bool, error)
	StoreInfo(user User, info string) error
	LoadInfo(user User) (string, error)
	PlayerNameMatchesId(user User) (bool, error)
}

type LibraryRepository interface {
	Store(library Library) (int, error)
	Remove(library Library) error
	FindById(id int) (Library, error, int)
}

type GameRepository interface {
	Store(game Game) (int, error)
	Remove(game Game) error
	FindById(id int) (Game, error, int)
	GameExisted(name string, libraryId int) (bool, error)
}

type User struct {
	Id           int
	Name         string
	Player       domain.Player //This user (account) was created by some player
	PersonalInfo string
	LibraryIds   []int
}

type Library struct {
	Id    int
	User  User //This library belongs to some user
	Games []Game
}

type Game struct {
	Id        int
	LibraryId int
	Name      string
	Producer  string
	Value     float64
}

type LoggerRepository interface {
	Log(message string) error
}

type ProfileInteractor struct {
	UserRepository    UserRepository
	LibraryRepository LibraryRepository
	GameRepository    GameRepository
	Loggr             LoggerRepository
}

func (interactor *ProfileInteractor) AddUser(player domain.Player, userName string) (int, error, int) {
	// Application rule: usernames cannot repeat
	existed, err := interactor.UserRepository.UserExisted(userName)
	if err != nil {
		return 0, err, 500
	}
	if existed {
		err := fmt.Errorf("Username '%s' is taken", userName)
		// interactor.Logger.Log(err.Error())
		return 0, err, 403
	}

	user := User{Name: userName, Player: player, PersonalInfo: ""}

	match, err := interactor.UserRepository.PlayerNameMatchesId(user)
	if err != nil {
		return 0, err, 500
	}
	if !match {
		err = fmt.Errorf("Player name does not match player Id")
		return 0, err, 400
	}

	id, err := interactor.UserRepository.Store(user)
	if err != nil {
		// interactor.Logger.Log(err.Error())
		return 0, err, 500
	}
	fmt.Printf("Added user #%d for player #%d\n", id, player.Id)
	return id, nil, 201
}

func (interactor *ProfileInteractor) ShowUser(userId int) (string, []int, error, int) {
	user, err, code := interactor.UserRepository.FindById(userId)
	if err != nil {
		err = fmt.Errorf(fmt.Sprintf("User #%d does not exist", userId))
		return "", nil, err, code
	}
	var libraryIds []int
	for _, libraryId := range user.LibraryIds {
		libraryIds = append(libraryIds, libraryId)
	}
	return user.Name, libraryIds, nil, 200
}

func (interactor *ProfileInteractor) RemoveUser(playerId, userId int) (error, int) {
	user, err, code := interactor.UserRepository.FindById(userId)
	if err != nil {
		// interactor.Logger.Log(err.Error())
		return err, code
	}

	if playerId != user.Player.Id {
		err := fmt.Errorf("Player #%d cannot remove user account of player #%d",
			playerId, user.Player.Id)
		// interactor.Logger.Log(err.Error())
		return err, 403
	}
	err = interactor.UserRepository.Remove(user)
	if err != nil {
		return err, 500
	}
	// interactor.Logger.Log(fmt.Sprintf("Removed user #%s (id #%d)", user.Name, user.Id))
	fmt.Printf("Player #%d deleted user #%d\n", playerId, user.Id)
	return nil, 200
}

func (interactor *ProfileInteractor) ShowUserInfo(userId int) (string, error, int) {
	user, err, code := interactor.UserRepository.FindById(userId)
	if err != nil {
		return "", err, code
	}
	info, err := interactor.UserRepository.LoadInfo(user)
	if err != nil {
		return "", err, 500
	}
	fmt.Println(fmt.Sprintf("Printed information of user #%d", user.Id))
	return info, nil, 200
}

func (interactor *ProfileInteractor) EditUserInfo(userId int, info string) (error, int) {
	user, err, code := interactor.UserRepository.FindById(userId)
	if err != nil {
		return err, code
	}
	err = interactor.UserRepository.StoreInfo(user, info)
	if err != nil {
		return err, 500
	}
	fmt.Println(fmt.Sprintf("Editted information of user '%s' (id #%d)", user.Name, user.Id))
	return nil, 200
}

func (interactor *ProfileInteractor) AddLibrary(userId int) (int, error, int) {
	user, err, code := interactor.UserRepository.FindById(userId)
	if err != nil {
		return 0, err, code
	}

	library := Library{User: user, Games: []Game{}}
	id, err := interactor.LibraryRepository.Store(library)
	if err != nil {
		return 0, err, 500
	}
	fmt.Printf("User #%d added library #%d\n", user.Id, id)
	return id, nil, 200
}

func (interactor *ProfileInteractor) ShowLibrary(userId, libraryId int) ([]Game, error, int) {
	var games []Game
	library, err, code := interactor.LibraryRepository.FindById(libraryId)
	if err != nil {
		err = fmt.Errorf(fmt.Sprintf("Library #%d of user #%d does not exist", libraryId, userId))
		return nil, err, code
	}

	if userId != library.User.Id {
		message := "User #%d is not allowed to see games in library #%d of user #%d"
		err := fmt.Errorf(message, userId, libraryId, library.User.Id)
		return nil, err, 403
	} else {
		for _, game := range library.Games {
			games = append(games, Game{Id: game.Id, Name: game.Name, Producer: game.Producer, Value: game.Value})
		}
		return games, nil, 200
	}
}

func (interactor *ProfileInteractor) RemoveLibrary(userId, libraryId int) (error, int) {
	user, err, code := interactor.UserRepository.FindById(userId)
	if err != nil {
		return err, code
	}
	library, err, code := interactor.LibraryRepository.FindById(libraryId)
	if err != nil {
		return err, code
	}
	if userId != library.User.Id {
		err := fmt.Errorf("User #%d cannot remove library of user #%d",
			userId, library.User.Id)
		return err, 403
	}

	for _, game := range library.Games {
		err = interactor.GameRepository.Remove(game)
		if err != nil {
			return err, 500
		}
	}
	err = interactor.LibraryRepository.Remove(library)
	if err != nil {
		return err, 500
	}
	fmt.Printf("User #%d removed library #%d\n", user.Id, library.Id)
	return nil, 200
}

func (interactor *ProfileInteractor) AddGame(userId, libraryId int, gameName, gameProducer string, gameValue float64) (int, error, int) {
	user, err, code := interactor.UserRepository.FindById(userId)
	if err != nil {
		// interactor.Logger.Log(err.Error())
		return 0, err, code
	}
	library, err, code := interactor.LibraryRepository.FindById(libraryId)
	if err != nil {
		// interactor.Logger.Log(err.Error())
		return 0, err, code
	}

	if user.Id != library.User.Id {
		message := "User #%d is not allowed to add games to library #%d of user #%d"
		err := fmt.Errorf(message, user.Id, library.Id, library.User.Id)
		// interactor.Logger.Log(err.Error())
		return 0, err, 403
	}
	existed, err := interactor.GameRepository.GameExisted(gameName, libraryId)
	if err != nil {
		return 0, err, 500
	}
	if existed {
		err := fmt.Errorf("Game '%s' is already in library", gameName)
		// interactor.Logger.Log(err.Error())
		return 0, err, 403
	}

	game := Game{LibraryId: library.Id, Name: gameName, Producer: gameProducer, Value: gameValue}
	id, err := interactor.GameRepository.Store(game)
	if err != nil {
		return 0, err, 500
	}

	fmt.Println(fmt.Sprintf("User added game %s (id #%d) to library #%d",
		game.Name, id, library.Id))
	return id, nil, 200
}

func (interactor *ProfileInteractor) RemoveGame(userId, libraryId, gameId int) (error, int) {
	user, err, code := interactor.UserRepository.FindById(userId)
	if err != nil {
		// interactor.Logger.Log(err.Error())
		return err, code
	}
	library, err, code := interactor.LibraryRepository.FindById(libraryId)
	if err != nil {
		// interactor.Logger.Log(err.Error())
		return err, code
	}
	if user.Player.Id != library.User.Player.Id {
		message := "User #%d is not allowed to remove games from library #%d of user #%d"
		err := fmt.Errorf(message, user.Id, library.Id, library.User.Id)
		// interactor.Logger.Log(err.Error())
		return err, 403
	}
	game, err, code := interactor.GameRepository.FindById(gameId)
	if err != nil {
		// interactor.Logger.Log(err.Error())
		return err, code
	}

	err = interactor.GameRepository.Remove(game)
	if err != nil {
		return err, 500
	}
	return nil, 200
}
