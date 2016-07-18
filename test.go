package main

import (
	"fmt"
	"game-tracker/domain"
	"game-tracker/usecases"
)

func main() {
	player := domain.Player{0, "Cam"}
	fmt.Println(player.Id)
	fmt.Println(player.Name)
}
