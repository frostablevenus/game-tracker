package interfaces

import (
	"database/sql"
	"fmt"

	"game-tracker/domain"
	"game-tracker/usecases"
)

type DbHandler interface {
	Execute(statement string, args ...interface{}) (sql.Result, error)
	Query(statement string, args ...interface{}) (Row, error)
}

type Row interface {
	Scan(dest ...interface{}) error
	Next() bool
	Close() error
}

type DbRepo struct {
	dbHandlers map[string]DbHandler
	dbHandler  DbHandler
}

type DbUserRepo DbRepo
type DbPlayerRepo DbRepo
type DbLibraryRepo DbRepo
type DbGameRepo DbRepo
type LoggerRepo DbRepo

func NewDbUserRepo(dbHandlers map[string]DbHandler) *DbUserRepo {
	dbUserRepo := new(DbUserRepo)
	dbUserRepo.dbHandlers = dbHandlers
	dbUserRepo.dbHandler = dbHandlers["DbUserRepo"]
	return dbUserRepo
}

func (repo DbUserRepo) Store(user usecases.User) error {
	_, err := repo.dbHandler.Execute(`INSERT INTO users (user_name, player_id, personal_info)
		VALUES ($1, $2, $3)`, user.Name, user.Player.Id, user.PersonalInfo)
	if err != nil {
		return err
	}
	playerRepo := NewDbPlayerRepo(repo.dbHandlers)
	err = playerRepo.Store(user.Player)
	return err
}

func (repo DbUserRepo) Remove(user usecases.User) error {
	_, err := repo.dbHandler.Execute(`DELETE FROM users WHERE id=$1`, user.Id)
	if err == sql.ErrNoRows {
	}
	return err
}

func (repo DbUserRepo) FindById(id int) (usecases.User, error) {
	row, err := repo.dbHandler.Query(`SELECT user_name, player_id, personal_info FROM users
		WHERE id = $1 LIMIT 1`, id)
	if err != nil {
		u := usecases.User{}
		return u, err
	}
	var userName string
	var playerId int
	var personalInfo string
	defer row.Close()
	row.Next()
	err = row.Scan(&userName, &playerId, &personalInfo)
	if err != nil {
		u := usecases.User{}
		return u, err
	}

	playerRepo := NewDbPlayerRepo(repo.dbHandlers)
	player, err := playerRepo.FindById(playerId)
	if err != nil {
		u := usecases.User{}
		return u, err
	}
	user := usecases.User{Id: id, Name: userName, Player: player, PersonalInfo: personalInfo}
	return user, nil
}

func (repo DbUserRepo) Count() (int, error) {
	row, err := repo.dbHandler.Query("SELECT COUNT(*) FROM users")
	if err != nil {
		return 0, err
	}
	var count int
	defer row.Close()
	row.Next()
	err = row.Scan(&count)
	return count, err
}

func (repo DbUserRepo) UserExisted(userName string) (bool, error) {
	row, err := repo.dbHandler.Query(`SELECT user_name FROM users
		WHERE user_name=$1 LIMIT 1`, userName)
	defer row.Close()
	return row.Next(), err
}

func (repo DbUserRepo) StoreInfo(user usecases.User, info string) error {
	_, err := repo.dbHandler.Execute(`UPDATE users SET personal_info=$1
		WHERE id=$2`, info, user.Id)
	return err
}

func (repo DbUserRepo) LoadInfo(user usecases.User) (string, error) {
	row, err := repo.dbHandler.Query(`SELECT personal_info FROM users WHERE id=$1`, user.Id)
	if err != nil {
		return "", err
	}
	var info string
	defer row.Close()
	row.Next()
	err = row.Scan(&info)
	return info, err
}

func NewDbPlayerRepo(dbHandlers map[string]DbHandler) *DbPlayerRepo {
	dbPlayerRepo := new(DbPlayerRepo)
	dbPlayerRepo.dbHandlers = dbHandlers
	dbPlayerRepo.dbHandler = dbHandlers["DbPlayerRepo"]
	return dbPlayerRepo
}

func (repo DbPlayerRepo) Store(player domain.Player) error {
	existed, err := repo.playerExisted(player.Name)
	if err != nil {
		return err
	}
	if !existed {
		_, err = repo.dbHandler.Execute(`INSERT INTO players (player_name)
		VALUES ($1)`, player.Name)
		return err
	}
	return nil
}

func (repo DbPlayerRepo) FindById(id int) (domain.Player, error) {
	row, err := repo.dbHandler.Query(`SELECT player_name FROM players WHERE id = $1 LIMIT 1`, id)
	if err != nil {
		p := domain.Player{}
		return p, err
	}
	var name string
	defer row.Close()
	row.Next()
	err = row.Scan(&name)
	if err != nil {
		p := domain.Player{}
		return p, err
	}
	return domain.Player{Id: id, Name: name}, nil
}

func (repo DbPlayerRepo) playerExisted(playerName string) (bool, error) {
	row, err := repo.dbHandler.Query(`SELECT player_name FROM players
		WHERE player_name=$1 LIMIT 1`, playerName)
	defer row.Close()
	return row.Next(), err
}

func NewDbLibraryRepo(dbHandlers map[string]DbHandler) *DbLibraryRepo {
	dbLibraryRepo := new(DbLibraryRepo)
	dbLibraryRepo.dbHandlers = dbHandlers
	dbLibraryRepo.dbHandler = dbHandlers["DbLibraryRepo"]
	return dbLibraryRepo
}

func (repo DbLibraryRepo) Store(library usecases.Library) error {
	if !repo.libraryExisted(library.Id) {
		_, err := repo.dbHandler.Execute(`INSERT INTO libraries (user_id) VALUES ($1)`,
			library.User.Id)
		if err != nil {
			return err
		}
	}

	for _, game := range library.Games {
		_, err := repo.dbHandler.Execute(`INSERT INTO games (game_id, library_id)
			VALUES ($1, $2)`, game.Id, library.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo DbLibraryRepo) Remove(library usecases.Library) error {
	_, err := repo.dbHandler.Execute(`DELETE FROM libraries WHERE id=$1`, library.Id)
	return err
}

func (repo DbLibraryRepo) FindById(id int) (usecases.Library, error) {
	row, err := repo.dbHandler.Query(`SELECT user_id FROM libraries WHERE id = $1 LIMIT 1`, id)
	if err != nil {
		return usecases.Library{}, err
	}

	var userId int
	defer row.Close()
	row.Next()
	row.Scan(&userId)
	userRepo := NewDbUserRepo(repo.dbHandlers)
	user, err := userRepo.FindById(userId)
	if err != nil {
		library := usecases.Library{}
		return library, err
	}
	library := usecases.Library{Id: id, User: user}

	var gameId int
	gameRepo := NewDbGameRepo(repo.dbHandlers)
	row, err = repo.dbHandler.Query(`SELECT game_id FROM games WHERE library_id = $1`, library.Id)
	if err != nil {
		return library, err
	}
	defer row.Close()
	for row.Next() {
		err = row.Scan(&gameId)
		if err != nil {
			return library, err
		}
		game, err := gameRepo.FindById(gameId)
		if err != nil {
			return library, err
		}
		library.Games = append(library.Games, game)
	}
	return library, err
}

func (repo DbLibraryRepo) libraryExisted(id int) bool {
	row, _ := repo.dbHandler.Query(`SELECT id FROM libraries
		WHERE id=$1 LIMIT 1`, id)
	defer row.Close()
	return row.Next()
}

func (repo DbLibraryRepo) Count() (int, error) {
	row, err := repo.dbHandler.Query(`SELECT COUNT(*) FROM libraries`)
	if err != nil {
		return 0, err
	}
	var count int
	defer row.Close()
	row.Next()
	err = row.Scan(&count)
	return count, err
}

func NewDbGameRepo(dbHandlers map[string]DbHandler) *DbGameRepo {
	dbGameRepo := new(DbGameRepo)
	dbGameRepo.dbHandlers = dbHandlers
	dbGameRepo.dbHandler = dbHandlers["DbGameRepo"]
	return dbGameRepo
}

func (repo DbGameRepo) Store(game usecases.Game) error {
	_, err := repo.dbHandler.Execute(`INSERT INTO games (library_id, game_name, producer, value)
    	VALUES ($1, $2, $3, $4)`, game.LibraryId, game.Name, game.Producer, game.Value)
	return err
}

func (repo DbGameRepo) Remove(game usecases.Game) error {
	_, err := repo.dbHandler.Execute(`DELETE FROM games WHERE game_id=$1`, game.Id+1)
	return err
}

func (repo DbGameRepo) FindById(id int) (usecases.Game, error) {
	row, err := repo.dbHandler.Query(`SELECT library_id, game_name, producer, value FROM games
    	WHERE game_id = $1 LIMIT 1`, id)
	if err != nil {
		game := usecases.Game{}
		return game, err
	}
	var (
		libraryId int
		name      string
		producer  string
		value     []uint8
	)

	defer row.Close()
	row.Next()
	err = row.Scan(&libraryId, &name, &producer, &value)
	if err != nil {
		return usecases.Game{}, err
	}

	game := usecases.Game{LibraryId: libraryId, Name: name, Producer: producer, Value: value}
	return game, nil
}

func (repo LoggerRepo) Log(message string) error {
	fmt.Println(message)
	return nil
}
