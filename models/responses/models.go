package responses

import (
	"fmt"
)

type Links struct {
	Self    string `json:"self,omitempty"`
	Related string `json:"related,omitempty"`
}

type Data struct {
	Type          string `json:"type,omitempty"`
	Id            int    `json:"id,omitempty"`
	Attributes    `json:"attributes,omitempty"`
	Relationships `json:"type, omitempty"`
}

type Attributes struct {
	TokenString string  `json:"tokenString,omitempty"`
	Name        string  `json:"name,omitempty"`
	Content     string  `json:"content,omitempty"`
	Producer    string  `json:"producer,omitempty"`
	Value       float64 `json:"value,omitempty"`
}

type Relationships struct {
	Libraries []Library `json:"libraries,omitempty"`
	Games     []Game    `json:"games,omitempty"`
	Owner     Owner     `json:"owner,omitempty"`
	Library   LibOfGame `json:"library,omitempty"`
}

type DataLv2 struct {
	Type       string `json:"type,omitempty"`
	Id         int    `json:"id,omitempty"`
	Attributes `json:"attributes,omitempty"`
}

type Token struct {
	Links `json:"links,omitempty"`
	Data  `json:"data, omitempty"`
}

type User struct {
	Links `json:"links,omitempty"`
	Data  `json:"data, omitempty"`
}

type Owner struct {
	DataLv2 `json:"data, omitempty"`
}

type Library struct {
	Links `json:"links,omitempty"`
	Data  `json:"data, omitempty"`
}

type LibOfGame struct {
	DataLv2
}

type Game struct {
	Links `json:"links,omitempty"`
	Data  `json:"data, omitempty"`
}

type Info struct {
	Links `json:"links,omitempty"`
	Data  `json:"data, omitempty"`
}

func ViewToken(tokenString string) Token {
	return Token{
		Data: Data{
			Type: "token",
			Attributes: Attributes{
				TokenString: tokenString,
			},
		},
	}
}

func ViewUser(id int, name string, libraries []Library) User {
	return User{
		Links: Links{
			Self: fmt.Sprintf("http://localhost:8080/users/%d", id),
		},
		Data: Data{
			Type: "users",
			Id:   id,
			Attributes: Attributes{
				Name: name,
			},
			Relationships: Relationships{
				Libraries: libraries,
			},
		},
	}
}

func ViewInfo(info string, userId int) Info {
	return Info{
		Links: Links{
			Self:    fmt.Sprintf("http://localhost:8080/users/%d/info", userId),
			Related: fmt.Sprintf("http://localhost:8080/users/%d", userId),
		},
		Data: Data{
			Type: "info",
			Id:   userId,
			Attributes: Attributes{
				Content: info,
			},
			Relationships: Relationships{
				Owner: Owner{
					DataLv2: DataLv2{
						Type: "users",
						Id:   userId,
					},
				},
			},
		},
	}
}

func ViewLibrary(userId, libId int, games []Game) Library {
	return Library{
		Links: Links{
			Self:    fmt.Sprintf("http://localhost:8080/users/%d/libraries/%d", userId, libId),
			Related: fmt.Sprintf("http://localhost:8080/users/%d", userId),
		},
		Data: Data{
			Type: "libraries",
			Id:   libId,
			Relationships: Relationships{
				Games: games,
				Owner: Owner{
					DataLv2: DataLv2{
						Type: "users",
						Id:   userId,
					},
				},
			},
		},
	}
}

func ViewGame(userId, libId, gameId int, name, producer string, value float64) Game {
	return Game{
		Links: Links{
			Self: fmt.Sprintf("http://localhost:8080/users/%d/libraries/%d/games/%d",
				userId, libId, gameId),
			Related: fmt.Sprintf("http://localhost:8080/users/%d/libraries/%d/",
				userId, libId),
		},
		Data: Data{
			Type: "games",
			Id:   gameId,
			Attributes: Attributes{
				Name:     name,
				Producer: producer,
				Value:    value,
			},
			Relationships: Relationships{
				Library: LibOfGame{
					DataLv2: DataLv2{
						Type: "libraries",
						Id:   libId,
					},
				},
			},
		},
	}
}

func ViewLibraries(libraryIds []int) []Library {
	var libraries []Library
	for _, id := range libraryIds {
		libraries = append(libraries, Library{
			Data: Data{
				Type: "libraries",
				Id:   id,
			},
		})
	}
	return libraries
}

func ViewGames(gameIds []int) []Game {
	var games []Game
	for _, id := range gameIds {
		games = append(games, Game{
			Data: Data{
				Type: "games",
				Id:   id,
			},
		})
	}
	return games
}
