package result

type User struct {
	Name       string
	Id         int
	LibraryIds []int
}

type UserInfo struct {
	Id   int
	Info string
}

type Game struct {
	Id       int
	Name     string
	Producer string
	Value    float64
}

type Library struct {
	Id     int
	UserId int
	Games  []Game
}
