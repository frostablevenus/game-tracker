package interfaces

import (
	"database/sql"

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
}

type DbRepo struct {
	dbHandlers map[string]DbHandler
	dbHandler  DbHandler
}

type DbUserRepo DbRepo
type DbPlayerRepo DbRepo
type DbLibraryRepo DbRepo
type DbGameRepo DbRepo

func NewDbUserRepo(dbHandlers map[string]DbHandler) *DbUserRepo {
	dbUserRepo := new(DbUserRepo)
	dbUserRepo.dbHandlers = dbHandlers
	dbUserRepo.dbHandler = dbHandlers["DbUserRepo"]
	return dbUserRepo
}

func (repo *DbUserRepo) Store(user usecases.User) error {
	_, err := repo.dbHandler.Execute(`INSERT INTO users (user_name, player_id, personal_info)
		VALUES ($1, $2, $3)`, user.Name, user.Player.Id, user.PersonalInfo)
	if err != nil {
		return err
	}
	playerRepo := NewDbPlayerRepo(repo.dbHandlers)
	err = playerRepo.Store(user.Player)
	return err
}

func (repo *DbUserRepo) Remove(user usecases.User) error {
	_, err := repo.dbHandler.Execute(`DELETE FROM users WHERE id=$1`, user.Id)
	return err
}

func (repo *DbUserRepo) FindById(id int) (usecases.User, error) {
	row, err := repo.dbHandler.Query(`SELECT (user_name, player_id, personal_info)
		FROM users WHERE id = $1 LIMIT 1`, id)
	if err != nil {
		u := usecases.User{}
		return u, err
	}
	var userName string
	var playerId int
	var personalInfo string
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

func (repo *DbUserRepo) Count() int {
	row, _ := repo.dbHandler.Query("SELECT user_name FROM users")
	count := 0
	for row.Next() {
		count++
	}
	return count
}

func (repo *DbUserRepo) NameExisted(userName string) bool {
	row, _ := repo.dbHandler.Query(`SELECT user_name FROM users
		WHERE user_name=$1 LIMIT 1`, userName)
	return row.Next()
}

func (repo *DbUserRepo) StoreInfo(user usecases.User, info string) error {
	_, err := repo.dbHandler.Execute(`UPDATE users SET personal_info=$1
		WHERE id=$2`, info, user.Id)
	return err
}

func (repo *DbUserRepo) LoadInfo(user usecases.User) (string, error) {
	row, err := repo.dbHandler.Query(`SELECT personal_info FROM users WHERE id=$1`, user.Id)
	if err != nil {
		return "", nil
	}
	var info string
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

func (repo *DbPlayerRepo) Store(player domain.Player) error {
	_, err := repo.dbHandler.Execute(`INSERT INTO players (name)
		VALUES ($1)`, player.Name)
	return err
}

func (repo *DbPlayerRepo) FindById(id int) (domain.Player, error) {
	row, err := repo.dbHandler.Query(`SELECT name FROM players WHERE id = $1 LIMIT 1`, id)
	if err != nil {
		p := domain.Player{}
		return p, err
	}
	var name string
	row.Next()
	err = row.Scan(&name)
	if err != nil {
		p := domain.Player{}
		return p, err
	}
	return domain.Player{Id: id, Name: name}, nil
}

func NewDbLibraryRepo(dbHandlers map[string]DbHandler) *DbLibraryRepo {
	dbLibraryRepo := new(DbLibraryRepo)
	dbLibraryRepo.dbHandlers = dbHandlers
	dbLibraryRepo.dbHandler = dbHandlers["DbLibraryRepo"]
	return dbLibraryRepo
}

func (repo *DbLibraryRepo) Store(library usecases.Library) error {
	_, err := repo.dbHandler.Execute(`INSERT INTO libraries (player_id) VALUES ($1)`,
		library.Player.Id)
	if err == nil {
		return err
	}
	for _, game := range library.Games {
		_, err = repo.dbHandler.Execute(`INSERT INTO gamesInLib (game_id, library_id)
			VALUES ($1, $2)`, game.Id, library.Id)
		if err == nil {
			return err
		}
	}
	return nil
}

func (repo *DbLibraryRepo) FindById(id int) (usecases.Library, error) {
	row, err := repo.dbHandler.Query(`SELECT player_id FROM libraries
		WHERE id = $1 LIMIT 1`, id)
	if err == nil {
		library := usecases.Library{}
		return library, err
	}

	var playerId int
	row.Next()
	row.Scan(&playerId)
	playerRepo := NewDbPlayerRepo(repo.dbHandlers)
	player, err := playerRepo.FindById(playerId)
	if err == nil {
		library := usecases.Library{}
		return library, err
	}
	library := usecases.Library{Id: id, Player: player}

	var gameId int
	gameRepo := NewDbGameRepo(repo.dbHandlers)
	row, err = repo.dbHandler.Query(`SELECT game_id FROM gamesInLib
		WHERE library_id = $1`, library.Id)
	if err == nil {
		return library, err
	}
	for row.Next() {
		err = row.Scan(&gameId)
		if err == nil {
			return library, err
		}
		game, err := gameRepo.FindById(gameId)
		if err == nil {
			return library, err
		}
		library.Games = append(library.Games, game)
	}
	return library, err
}

func NewDbGameRepo(dbHandlers map[string]DbHandler) *DbGameRepo {
	dbGameRepo := new(DbGameRepo)
	dbGameRepo.dbHandlers = dbHandlers
	dbGameRepo.dbHandler = dbHandlers["DbGameRepo"]
	return dbGameRepo
}

func (repo *DbGameRepo) Store(game domain.Game) error {
	_, err := repo.dbHandler.Execute(`INSERT INTO games (game_name, producer, value)
    	VALUES ($1, $2, $3)`, game.Name, game.Producer, game.Value)
	if err == nil {
		return err
	}
	return nil
}

func (repo *DbGameRepo) FindById(id int) (domain.Game, error) {
	row, err := repo.dbHandler.Query(`SELECT game_name, producer, value FROM items
    	WHERE id = $1 LIMIT 1`, id)
	if err == nil {
		game := domain.Game{}
		return game, err
	}
	var name string
	var producer string
	var value float64
	row.Next()
	err = row.Scan(&name, &producer, &value)
	if err == nil {
		return domain.Game{}, err
	}
	game := domain.Game{Id: id, Name: name, Producer: producer, Value: value}
	return game, nil
}
