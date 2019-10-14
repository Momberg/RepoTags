package main

import (
	"context"
	"fmt"
	"log"
	"repotags/model"
	"repotags/repository"

	"github.com/google/go-github/github"
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
}

func main() {
	/*gitRepo := model.Repo{
		Description: "teste",
		ID:          "teste",
		Language:    "EN",
		URL:         "www.google.com.br",
	}
	gitRepo.Tags = append(gitRepo.Tags, "a")
	gitRepo.Tags = append(gitRepo.Tags, "b")
	if model.ValidateDuplicatedTag(gitRepo, "a") {
		fmt.Println("This tag already exists")
	}
	fmt.Printf("git %+v\r\n ", gitRepo)*/

	//http.HandleFunc("/", handler.CreateRestHandler)
	//http.ListenAndServe(":8181", nil)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "add your token"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	repos, _, err := client.Activity.ListStarred(ctx, "", nil)
	if err != nil {
		fmt.Println("Error to access gitHub: ", err.Error())
	}
	var gitRepo model.Repository
	db, err := repository.GetDBConnection()
	if err != nil {
		log.Println("[Main] Connection error: ", err.Error())
		return
	}
	sql := "insert into repository (id, name, description, url, language) values (?, ?, ?, ?, ?)"
	for _, starredRepo := range repos {
		gitRepo = model.Repository{
			Description: starredRepo.GetRepository().GetDescription(),
			ID:          starredRepo.GetRepository().GetID(),
			Language:    starredRepo.GetRepository().GetLanguage(),
			URL:         starredRepo.GetRepository().GetURL(),
			Name:        starredRepo.Repository.GetName(),
		}
		fmt.Printf("User repos: %+v\r\n", gitRepo)
		_, err := db.Exec(sql, starredRepo.GetRepository().GetID(), starredRepo.Repository.GetName(), starredRepo.GetRepository().GetDescription(),
			starredRepo.GetRepository().GetURL(), starredRepo.GetRepository().GetLanguage())
		if err != nil {
			fmt.Println("[Main] Insert error: ", sql, " - ", err.Error())
		}
	}
}
