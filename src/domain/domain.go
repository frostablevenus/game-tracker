package domain

type PlayerRepository interface {
	Store(player Player)
	FindById(id int) Player
}

type GameRepository interface {
	Store(game Game)
	FindById(id int) Game
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
