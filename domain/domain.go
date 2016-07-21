package domain

type PlayerRepository interface {
	Store(player Player) error
	FindById(id int) (Player, error, int)
	NameMatchesId(playerName string, id int) (bool, error)
}

type Player struct {
	Id   int
	Name string
}

type Game struct {
	Id       int
	Name     string
	Producer string
	Value    []uint8
}

//Business rule: Player names cannot repeat (unique identification)
