package domain

type PlayerRepository interface {
	Store(player Player) error
	FindById(id int) (Player, error)
}

type GameRepository interface {
	Store(game Game) error
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
