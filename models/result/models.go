package result

type User struct {
	Id         int    `json:"UserId"`
	Name       string `json:"name"`
	LibraryIds []int  `json:"libraryIds"`
}

type UserAdd struct {
	Id   int    `json:"userId"`
	Name string `json:"name"`
}

type UserDelete struct {
	Id int `json:"userId"`
}

type UserInfo struct {
	Id   int    `json:"userId"`
	Info string `json:"userInfo"`
}

type Game struct {
	Id        int     `json:"gameId"`
	LibraryId int     `json:"libraryId"`
	UserId    int     `json:"userId"`
	Name      string  `json:"name"`
	Producer  string  `json:"producer"`
	Value     float64 `json:"value"`
}

type GameToLib struct {
	Id        int `json:"gameId"`
	LibraryId int `json:"libraryId"`
	UserId    int `json:"userId"`
}

type Library struct {
	Id       int   `json:"libraryId"`
	UserId   int   `json:"userId"`
	GamesIds []int `json:"gameIds"`
}

type LibraryAdd struct {
	Id     int `json:"libraryId"`
	UserId int `json:"userId"`
}

type LibraryDelete struct {
	Id int `json:"libraryId"`
}
