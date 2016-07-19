package usecases

import (
	"fmt"

	"game-tracker/domain"
)

type UserRepository interface {
	Store(user User) (int, error)
	Remove(user User) error
	FindById(id int) (User, error)
	UserExisted(userName string) (bool, error)
	StoreInfo(user User, info string) error
	LoadInfo(user User) (string, error)
}

type LibraryRepository interface {
	Store(library Library) (int, error)
	Remove(library Library) error
	FindById(id int) (Library, error)
}

type GameRepository interface {
	Store(game Game) (int, error)
	Remove(game Game) error
	FindById(id int) (Game, error)
}

type User struct {
	Id           int
	Name         string
	Player       domain.Player //This user (account) was created by some player
	PersonalInfo string
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
	Value     []uint8
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

func (interactor *ProfileInteractor) ShowLibrary(userId, libraryId int) (string, []Game, error) {
	var games []Game
	user, err := interactor.UserRepository.FindById(userId)
	if err != nil {
		err = fmt.Errorf(fmt.Sprintf("User #%d does not exist", userId))
		return "", nil, err
	}
	info, err := interactor.UserRepository.LoadInfo(user)
	if err != nil {
		err = fmt.Errorf("Error loading user personal information")
		return "", nil, err
	}
	library, err := interactor.LibraryRepository.FindById(libraryId)
	if err != nil {
		err = fmt.Errorf(fmt.Sprintf("Library #%d of user #%d does not exist", libraryId, userId))
		return "", nil, err
	}

	if user.Id != libraryId {
		message := "User #%d is not allowed to see games in library #%d of user #%d"
		err := fmt.Errorf(message, user.Id, libraryId, library.User.Id)
		games = make([]Game, 0)
		return info, games, err
	} else {
		games = make([]Game, len(library.Games))
		for i, game := range library.Games {
			games[i] = Game{Id: game.Id, Name: game.Name, Producer: game.Producer, Value: game.Value}
		}
		return info, games, nil
	}
}

func (interactor *ProfileInteractor) AddUser(player domain.Player, userName string) (int, error) {
	// Application rule: usernames cannot repeat
	existed, err := interactor.UserRepository.UserExisted(userName)
	if err != nil {
		return 0, err
	}
	if existed {
		err := fmt.Errorf("Username '%s' is taken", userName)
		// interactor.Logger.Log(err.Error())
		return 0, err
	}
	user := User{Name: userName, Player: player, PersonalInfo: ""}
	id, err := interactor.UserRepository.Store(user)
	if err != nil {
		// interactor.Logger.Log(err.Error())
		return 0, err
	}
	fmt.Printf("Added user #%d for player #%d\n", id, player.Id)
	return id, nil
}

func (interactor *ProfileInteractor) RemoveUser(playerId, userId int) error {
	user, err := interactor.UserRepository.FindById(userId)
	if err != nil {
		// interactor.Logger.Log(err.Error())
		return err
	}

	if playerId != user.Player.Id {
		err := fmt.Errorf("Player #%d cannot remove user account of player #%d",
			playerId, user.Player.Id)
		// interactor.Logger.Log(err.Error())
		return err
	}
	interactor.UserRepository.Remove(user)
	// interactor.Logger.Log(fmt.Sprintf("Removed user #%s (id #%d)", user.Name, user.Id))
	fmt.Printf("Player #%d deleted user #%d\n", playerId, user.Id)
	return nil
}

func (interactor *ProfileInteractor) EditUserInfo(userId, targetId int, info string) error {
	if userId != targetId {
		err := fmt.Errorf("User #%d cannot modify information of user #%d", userId, targetId)
		return err
	}
	user, err := interactor.UserRepository.FindById(userId)
	if err != nil {
		return err
	}
	err = interactor.UserRepository.StoreInfo(user, info)
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("Editted information of user '%s' (id #%d)", user.Name, user.Id))
	return nil
}

func (interactor *ProfileInteractor) AddLibrary(userId int) (int, error) {
	user, err := interactor.UserRepository.FindById(userId)
	if err != nil {
		return 0, err
	}

	library := Library{User: user, Games: []Game{}}
	id, err := interactor.LibraryRepository.Store(library)
	if err != nil {
		return 0, err
	}
	fmt.Printf("User #%d added library #%d\n", user.Id, id)
	return id, nil
}

func (interactor *ProfileInteractor) RemoveLibrary(userId, libraryId int) error {
	user, err := interactor.UserRepository.FindById(userId)
	if err != nil {
		return err
	}
	library, err := interactor.LibraryRepository.FindById(libraryId)
	if err != nil {
		return err
	}
	if userId != library.User.Id {
		err := fmt.Errorf("User #%d cannot remove library of user #%d",
			userId, library.User.Id)
		return err
	}

	for _, game := range library.Games {
		err = interactor.GameRepository.Remove(game)
		if err != nil {
			return err
		}
	}
	err = interactor.LibraryRepository.Remove(library)
	if err != nil {
		return err
	}
	fmt.Printf("User #%d removed library #%d\n", user.Id, library.Id)
	return nil
}

func (interactor *ProfileInteractor) AddGame(userId, libraryId int, gameName, gameProducer string, gameValue []uint8) (int, error) {
	user, err := interactor.UserRepository.FindById(userId)
	if err != nil {
		// interactor.Logger.Log(err.Error())
		return 0, err
	}
	library, err := interactor.LibraryRepository.FindById(libraryId)
	if err != nil {
		// interactor.Logger.Log(err.Error())
		return 0, err
	}

	if user.Id != library.User.Id {
		message := "User #%d is not allowed to add games to library #%d of user #%d"
		err := fmt.Errorf(message, user.Id, library.Id, library.User.Id)
		// interactor.Logger.Log(err.Error())
		return 0, err
	}

	game := Game{LibraryId: library.Id, Name: gameName, Producer: gameProducer, Value: gameValue}
	id, err := interactor.GameRepository.Store(game)
	if err != nil {
		return 0, err
	}

	fmt.Println(fmt.Sprintf("User added game %s (id #%d) to library #%d",
		game.Name, id, library.Id))
	return id, nil
}

func (interactor *ProfileInteractor) RemoveGame(userId, libraryId, gameId int) error {
	user, err := interactor.UserRepository.FindById(userId)
	if err != nil {
		// interactor.Logger.Log(err.Error())
		return err
	}
	library, err := interactor.LibraryRepository.FindById(libraryId)
	if err != nil {
		// interactor.Logger.Log(err.Error())
		return err
	}
	if user.Player.Id != library.User.Player.Id {
		message := "User #%d is not allowed to remove games from library #%d of user #%d"
		err := fmt.Errorf(message, user.Id, library.Id, library.User.Id)
		// interactor.Logger.Log(err.Error())
		return err
	}
	game, err := interactor.GameRepository.FindById(gameId)
	if err != nil {
		// interactor.Logger.Log(err.Error())
		return err
	}

	err = interactor.GameRepository.Remove(game)
	if err != nil {
		return err
	}
	return nil
}
