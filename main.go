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

	http.HandleFunc("/show_library", func(res http.ResponseWriter, req *http.Request) {
		err := webserviceHandler.ShowLibrary(res, req)
		if err != nil {
			fmt.Println(err)
		}
	})

	http.HandleFunc("/add_user", func(res http.ResponseWriter, req *http.Request) {
		err := webserviceHandler.AddUser(res, req)
		if err != nil {
			fmt.Println(err)
		}
	})

	http.HandleFunc("/remove_user", func(res http.ResponseWriter, req *http.Request) {
		err := webserviceHandler.RemoveUser(res, req)
		if err != nil {
			fmt.Println(err)
		}
	})

	http.HandleFunc("/edit_user_info", func(res http.ResponseWriter, req *http.Request) {
		err := webserviceHandler.EditUserInfo(res, req)
		if err != nil {
			fmt.Println(err)
		}
	})

	http.HandleFunc("/add_library", func(res http.ResponseWriter, req *http.Request) {
		err := webserviceHandler.AddLibrary(res, req)
		if err != nil {
			fmt.Println(err)
		}
	})

	http.HandleFunc("/remove_library", func(res http.ResponseWriter, req *http.Request) {
		err := webserviceHandler.RemoveLibrary(res, req)
		if err != nil {
			fmt.Println(err)
		}
	})

	http.HandleFunc("/add_game", func(res http.ResponseWriter, req *http.Request) {
		err := webserviceHandler.AddGame(res, req)
		if err != nil {
			fmt.Println(err)
		}
	})

	http.HandleFunc("/remove_game", func(res http.ResponseWriter, req *http.Request) {
		err := webserviceHandler.RemoveGame(res, req)
		if err != nil {
			fmt.Println(err)
		}
	})

	fmt.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}
