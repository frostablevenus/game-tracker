package main

import (
	"encoding/json"
	"fmt"
	"os"

	"game-tracker/infrastructure"
	"game-tracker/interfaces"
	"game-tracker/models/postgres"
	"game-tracker/routes"
	"game-tracker/usecases"
)

func main() {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Cannot open config file")
		return
	}
	decoder := json.NewDecoder(file)
	config := postgres.Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Cannot read config file")
		return
	}
	dbHandler, err := infrastructure.NewPostgresqlHandler(config.PostgresAdr)
	if err != nil {
		fmt.Println("Cannot open database", err)
		return
	}

	handlers := make(map[string]interfaces.DbHandler)
	handlers["DbUserRepo"] = dbHandler
	handlers["DbPlayerRepo"] = dbHandler
	handlers["DbGameRepo"] = dbHandler
	handlers["DbLibraryRepo"] = dbHandler

	profileInteractor := usecases.ProfileInteractor{
		UserRepository:    interfaces.NewDbUserRepo(handlers),
		GameRepository:    interfaces.NewDbGameRepo(handlers),
		LibraryRepository: interfaces.NewDbLibraryRepo(handlers),
	}

	webserviceHandler := interfaces.WebserviceHandler{}
	webserviceHandler.ProfileInteractor = profileInteractor

	engine := routes.CreateEngine(webserviceHandler)

	fmt.Println("Listening...")
	engine.Run(":8080")
}
