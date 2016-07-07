package domain

type PlayerRepository interface {
	Store(player Player)
	FindById(id int) (Player, error)
}

type GameRepository interface {
	Store(game Game)
	FindById(id int) (Game, error)
}

type Player struct {
	Id   int
	Name string
}

type Game struct {
	Id       int
	Name     string
	Producer string
	Value    float64
}
