package request

type UserInfo struct {
	Info string `json:"info" binding:"required"`
}

type Game struct {
	Name     string  `json:"name" binding:"required"`
	Producer string  `json:"producer" binding:"required"`
	Value    float64 `json:"value" binding:"required"`
}

type User struct {
	PlayerId   int    `json:"playerId" binding:"required"`
	PlayerName string `json:"playerName" binding:"required"`
	Name       string `json:"name" binding:"required"`
}
