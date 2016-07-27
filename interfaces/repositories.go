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
	QueryRow(statement string, args ...interface{}) (int, error)
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

func (repo DbUserRepo) Store(user usecases.User) (int, error) {
	id, err := repo.dbHandler.QueryRow(`INSERT INTO users (user_name, player_id, personal_info)
		VALUES ($1, $2, $3) RETURNING id`, user.Name, user.Player.Id, user.PersonalInfo)
	if err != nil {
		return 0, err
	}

	playerRepo := NewDbPlayerRepo(repo.dbHandlers)
	err = playerRepo.Store(user.Player)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (repo DbUserRepo) Remove(user usecases.User) error {
	_, err := repo.dbHandler.Execute(`DELETE FROM users WHERE id=$1`, user.Id)
	return err
}

func (repo DbUserRepo) FindById(id int) (usecases.User, error, int) {
	row, err := repo.dbHandler.Query(`SELECT user_name, player_id, personal_info FROM users
		WHERE id = $1 LIMIT 1`, id)
	if err != nil {
		return usecases.User{}, err, 500
	}
	var userName string
	var playerId int
	var personalInfo string
	defer row.Close()
	row.Next()
	err = row.Scan(&userName, &playerId, &personalInfo)
	if err != nil {
		return usecases.User{}, err, 404
	}

	playerRepo := NewDbPlayerRepo(repo.dbHandlers)
	player, err, code := playerRepo.FindById(playerId)
	if err != nil {
		return usecases.User{}, err, code
	}

	user := usecases.User{Id: id, Name: userName, Player: player, PersonalInfo: personalInfo}

	var libraryId int
	row, err = repo.dbHandler.Query(`SELECT id FROM libraries WHERE user_id = $1`, id)
	if err != nil {
		return user, err, 500
	}
	defer row.Close()
	for row.Next() {
		err = row.Scan(&libraryId)
		if err != nil {
			return user, err, 500
		}
		user.LibraryIds = append(user.LibraryIds, libraryId)
	}
	return user, nil, 200
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

func (repo DbUserRepo) PlayerNameMatchesId(user usecases.User) (bool, error) {
	playerRepo := NewDbPlayerRepo(repo.dbHandlers)
	existed, err := playerRepo.playerExisted(user.Player.Name)
	if err != nil {
		return false, err
	}
	if existed {
		match, err := playerRepo.nameMatchesId(user.Player.Name, user.Player.Id)
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}
	return true, nil
}

func (repo DbUserRepo) AddLoginInfo(username, password string) error {
	_, err := repo.dbHandler.Execute(`INSERT INTO loginInfo (username, password)
		VALUES ($1, $2)`, username, password)
	if err != nil {
		return err
	}
	return nil
}

func (repo DbUserRepo) FindLoginId(username, password string) (int, bool, error) {
	row, err := repo.dbHandler.Query(`SELECT id FROM loginInfo WHERE username=$1
		AND password=$2 LIMIT 1`, username, password)
	if err != nil {
		return 0, false, err
	}

	var id int
	defer row.Close()
	exist := row.Next()
	if !exist {
		return 0, false, nil
	}
	err = row.Scan(&id)
	if err != nil {
		return 0, true, err
	}
	return id, true, nil
}

func (repo DbUserRepo) RemoveLoginInfo(user usecases.User) error {
	_, err := repo.dbHandler.Execute(`DELETE FROM loginInfo WHERE username=$1`, user.Name)
	return err
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
	if existed {
		return nil
	}
	_, err = repo.dbHandler.Execute(`INSERT INTO players (player_name)
		VALUES ($1)`, player.Name)
	return err
}

func (repo DbPlayerRepo) FindById(id int) (domain.Player, error, int) {
	row, err := repo.dbHandler.Query(`SELECT player_name FROM players WHERE id = $1 LIMIT 1`, id)
	if err != nil {
		return domain.Player{}, err, 500
	}
	var name string
	defer row.Close()
	row.Next()
	err = row.Scan(&name)
	if err != nil {
		return domain.Player{}, err, 404
	}
	return domain.Player{Id: id, Name: name}, nil, 200
}

func (repo DbPlayerRepo) playerExisted(playerName string) (bool, error) {
	row, err := repo.dbHandler.Query(`SELECT player_name FROM players
		WHERE player_name=$1 LIMIT 1`, playerName)
	defer row.Close()
	return row.Next(), err
}

func (repo DbPlayerRepo) nameMatchesId(playerName string, id int) (bool, error) {
	row, err := repo.dbHandler.Query(`SELECT * FROM players
		WHERE id=$1 AND player_name=$2 LIMIT 1`, id, playerName)
	defer row.Close()
	return row.Next(), err
}

func NewDbLibraryRepo(dbHandlers map[string]DbHandler) *DbLibraryRepo {
	dbLibraryRepo := new(DbLibraryRepo)
	dbLibraryRepo.dbHandlers = dbHandlers
	dbLibraryRepo.dbHandler = dbHandlers["DbLibraryRepo"]
	return dbLibraryRepo
}

func (repo DbLibraryRepo) Store(library usecases.Library) (int, error) {
	id, err := repo.dbHandler.QueryRow(`INSERT INTO libraries (user_id) VALUES ($1) RETURNING id`,
		library.User.Id)
	return id, err
}

func (repo DbLibraryRepo) Remove(library usecases.Library) error {
	_, err := repo.dbHandler.Execute(`DELETE FROM libraries WHERE id=$1`, library.Id)
	return err
}

func (repo DbLibraryRepo) FindById(id int) (usecases.Library, error, int) {
	row, err := repo.dbHandler.Query(`SELECT user_id FROM libraries WHERE id = $1 LIMIT 1`, id)
	if err != nil {
		return usecases.Library{}, err, 500
	}

	var userId int
	defer row.Close()
	row.Next()
	err = row.Scan(&userId)
	if err != nil {
		return usecases.Library{}, err, 404
	}
	userRepo := NewDbUserRepo(repo.dbHandlers)
	user, err, code := userRepo.FindById(userId)
	if err != nil {
		return usecases.Library{}, err, code
	}
	library := usecases.Library{Id: id, User: user}

	var gameId int
	row, err = repo.dbHandler.Query(`SELECT id FROM gamesInLib WHERE library_id = $1`, library.Id)
	if err != nil {
		return library, err, 500
	}
	defer row.Close()
	for row.Next() {
		err = row.Scan(&gameId)
		if err != nil {
			return library, err, 404
		}
		library.GameIds = append(library.GameIds, gameId)
	}
	return library, err, 200
}

func NewDbGameRepo(dbHandlers map[string]DbHandler) *DbGameRepo {
	dbGameRepo := new(DbGameRepo)
	dbGameRepo.dbHandlers = dbHandlers
	dbGameRepo.dbHandler = dbHandlers["DbGameRepo"]
	return dbGameRepo
}

func (repo DbGameRepo) Store(game usecases.Game) (int, error) {
	id, existed, err := repo.gameExisted(game.Name)
	if !existed {
		id, err = repo.dbHandler.QueryRow(`INSERT INTO games (name, producer, value)
    	VALUES ($1, $2, $3) RETURNING id`, game.Name, game.Producer, game.Value)
		return id, err
	}
	return id, nil
}

func (repo DbGameRepo) AddToLib(gameId, libraryId int) (error, int) {
	existed, err := repo.gameExistedInLib(gameId, libraryId)
	if err != nil {
		return err, 500
	}
	if existed {
		err = fmt.Errorf("Game already existed in library")
		return err, 400
	}
	_, err = repo.dbHandler.Execute(`INSERT INTO gamesInLib (game_id, library_id)
		VALUES ($1, $2)`, gameId, libraryId)
	if err != nil {
		return err, 500
	}
	return nil, 200
}

func (repo DbGameRepo) RemoveFromLib(game usecases.Game, libraryId int) error {
	_, err := repo.dbHandler.Execute(`DELETE FROM gamesInLib WHERE game_id=$1 AND library_id=$2`,
		game.Id, libraryId)
	return err
}

func (repo DbGameRepo) gameExisted(name string) (int, bool, error) {
	row, err := repo.dbHandler.Query(`SELECT id FROM games
		WHERE name=$1 LIMIT 1`, name)
	if err != nil {
		return 0, false, err
	}

	exist := row.Next()
	if !exist {
		return 0, false, nil
	}
	var id int
	defer row.Close()
	err = row.Scan(&id)
	if err != nil {
		return 0, false, err
	}
	return id, true, nil
}

func (repo DbGameRepo) gameExistedInLib(gameId, libraryId int) (bool, error) {
	row, err := repo.dbHandler.Query(`SELECT id FROM gamesInLib
		WHERE game_id=$1 AND library_id=$2 LIMIT 1`, gameId, libraryId)
	if err != nil {
		return false, err
	}
	defer row.Close()
	return row.Next(), nil
}

func (repo DbGameRepo) FindById(id int) (usecases.Game, error, int) {
	row, err := repo.dbHandler.Query(`SELECT name, producer, value FROM games
    	WHERE id = $1 LIMIT 1`, id)
	if err != nil {
		return usecases.Game{}, err, 500
	}
	var (
		name     string
		producer string
		value    float64
	)

	defer row.Close()
	row.Next()
	err = row.Scan(&name, &producer, &value)
	if err != nil {
		return usecases.Game{}, err, 404
	}

	game := usecases.Game{Id: id, Name: name, Producer: producer, Value: value}
	return game, nil, 200
}

func (repo LoggerRepo) Log(message string) error {
	fmt.Println(message)
	return nil
}
