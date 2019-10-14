package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"repotags/handler"
	"repotags/repository"

	"github.com/go-sql-driver/mysql"
	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

func init() {
	fmt.Println("Starting DB connection...")
	_, err := repository.OpenConnection()
	if err != nil {
		fmt.Println("Stop loading. Error on DB oppening: ", err.Error())
		return
	}
	fmt.Println("DB connection started.")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "f59811e8859a70925ae34fa5eef6cda09ae016aa"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	repos, _, err := client.Activity.ListStarred(ctx, "", nil)
	if err != nil {
		fmt.Println("Error to access gitHub: ", err.Error())
	}
	db, err := repository.GetDBConnection()
	if err != nil {
		log.Println("[Main] DB connection error: ", err.Error())
		return
	}
	sql := "insert into repository (id, name, description, url, language) values (?, ?, ?, ?, ?)"
	for _, starredRepo := range repos {
		_, err := db.Exec(sql, starredRepo.GetRepository().GetID(), starredRepo.Repository.GetName(), starredRepo.GetRepository().GetDescription(),
			starredRepo.GetRepository().GetURL(), starredRepo.GetRepository().GetLanguage())
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == 1062 {
				return
			}
			fmt.Println("[Main] Insert error: ", sql, " - ", driverErr.Error())
		}

	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/repositories", handler.GetRepos).Methods("GET")
	router.HandleFunc("/repository/{id}/tag", handler.AddTagToRepo).Methods("POST")
	http.ListenAndServe(":8181", router)
}
