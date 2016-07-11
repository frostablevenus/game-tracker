package main

import (
	"fmt"
	"net/http"

	"game-tracker/infrastructure"
	"game-tracker/interfaces"
	"game-tracker/usecases"
)

func main() {
	dbAddress := "postgres://postgres:postgres5263@127.0.0.1:5432/gameTrackerDb?sslmode=disable"
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
