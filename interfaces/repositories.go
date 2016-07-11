package interfaces

import (
	"fmt"

	"game-tracker/domain"
	"game-tracker/usecases"
)

type DbHandler interface {
	Execute(statement string) error
	Query(statement string) (Row, error)
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
	err := repo.dbHandler.Execute(fmt.Sprintf(`INSERT INTO users (id, user_name, player_id, personal_info)
		VALUES ('%d', '%s', '%d', '%s')`,
		user.Id, user.Name, user.Player.Id, user.PersonalInfo))
	if err != nil {
		return err
	}
	playerRepo := NewDbPlayerRepo(repo.dbHandlers)
	err = playerRepo.Store(user.Player)
	return err
}

func (repo *DbUserRepo) FindById(id int) (usecases.User, error) {
	row, err := repo.dbHandler.Query(fmt.Sprintf(`SELECT (user_name, player_id, personal_info)
		FROM users WHERE id = '%d' LIMIT 1`, id))
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

func NewDbPlayerRepo(dbHandlers map[string]DbHandler) *DbPlayerRepo {
	dbPlayerRepo := new(DbPlayerRepo)
	dbPlayerRepo.dbHandlers = dbHandlers
	dbPlayerRepo.dbHandler = dbHandlers["DbPlayerRepo"]
	return dbPlayerRepo
}

func (repo *DbPlayerRepo) Store(player domain.Player) error {
	err := repo.dbHandler.Execute(fmt.Sprintf(`INSERT INTO players (id, name)
		VALUES ('%d', '%v')`,
		player.Id, player.Name))
	return err
}

func (repo *DbPlayerRepo) FindById(id int) (domain.Player, error) {
	row, err := repo.dbHandler.Query(fmt.Sprintf(`SELECT name FROM players
		WHERE id = '%d' LIMIT 1`, id))
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
	err := repo.dbHandler.Execute(fmt.Sprintf(`INSERT INTO libraries (id, player_id) VALUES ('%d', '%v')`,
		library.Id, library.Player.Id))
	if err == nil {
		return err
	}
	for _, item := range library.Games {
		err = repo.dbHandler.Execute(fmt.Sprintf(`INSERT INTO gamesInLib (game_id, library_id)
			VALUES ('%d', '%d')`, item.Id, library.Id))
		if err == nil {
			return err
		}
	}
	return nil
}

func (repo *DbLibraryRepo) FindById(id int) (usecases.Library, error) {
	row, err := repo.dbHandler.Query(fmt.Sprintf(`SELECT player_id FROM libraries
		WHERE id = '%d' LIMIT 1`, id))
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
	row, err = repo.dbHandler.Query(fmt.Sprintf(`SELECT game_id FROM gamesInLib
		WHERE library_id = '%d'`, library.Id))
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
	err := repo.dbHandler.Execute(fmt.Sprintf(`INSERT INTO games (id, name, producer value)
    	VALUES ('%d', '%s', '%s', '%f')`, game.Id, game.Name, game.Producer, game.Value))
	if err == nil {
		return err
	}
	return nil
}

func (repo *DbGameRepo) FindById(id int) (domain.Game, error) {
	row, err := repo.dbHandler.Query(fmt.Sprintf(`SELECT name, producer, value
    	FROM items WHERE id = '%d' LIMIT 1`, id))
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
