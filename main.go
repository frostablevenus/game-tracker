package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"game-tracker/infrastructure"
	"game-tracker/interfaces"
	"game-tracker/models"
	"game-tracker/usecases"
)

func main() {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Cannot open config file")
		return
	}
	decoder := json.NewDecoder(file)
	config := models.Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Cannot read config file")
		return
	}
	dbAddress := config.PostgresAdr
	dbHandler, err := infrastructure.NewPostgresqlHandler(dbAddress)
	if err != nil {
		fmt.Println("Cannot open database")
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

	http.HandleFunc("/library", func(res http.ResponseWriter, req *http.Request) {
		webserviceHandler.ShowLibrary(res, req)
	})
	http.ListenAndServe(":8080", nil)
}
