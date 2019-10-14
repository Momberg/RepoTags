package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"repotags/model"
	"repotags/repository"

	"github.com/gorilla/mux"
)

//GetRepos get repositories
func GetRepos(w http.ResponseWriter, r *http.Request) {
	db, err := repository.GetDBConnection()
	if err != nil {
		log.Println("[Rest] Connection error: ", err.Error())
		return
	}
	sql := "select id, name, description, url, language from repository"
	repo := []model.Repository{}
	err = db.Select(&repo, sql)
	if err != nil {
		log.Println("[Rest] Select error: ", err.Error())
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repo)
}

//AddTagToRepo add tag to repository
func AddTagToRepo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if len(params["id"]) < 1 {
		log.Println("[Rest] Missing repository ID.")
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	tag := model.Tags{}
	err = json.Unmarshal(body, &tag)
	if err != nil {
		log.Println("[Rest] Body convertion error: ", err.Error())
		return
	}
	if validateTags(tag.Name) {
		http.Error(w, "This tag already exists for this repository.", http.StatusBadRequest)
		return
	}
	db, err := repository.GetDBConnection()
	if err != nil {
		log.Println("[Rest] Connection error: ", err.Error())
		return
	}
	sql := "insert into tags (id, name) values (?, ?)"
	_, err = db.Exec(sql, params["id"], &tag.Name)
	if err != nil {
		log.Println("[Rest] Insert tags error: ", err.Error())
	}
	sql = "update repository set tag_id = (?) where id = " + params["id"]
	_, err = db.Exec(sql, params["id"])
	if err != nil {
		log.Println("[Rest] Update repository error: ", err.Error())
	}
}

//ValidateTags validate tags
func validateTags(name string) bool {
	db, err := repository.GetDBConnection()
	if err != nil {
		log.Println("[Rest] Connection error: ", err.Error())
	}
	sql := "select name from tags"
	tags := []model.Tags{}
	err = db.Select(&tags, sql)
	if err != nil {
		log.Println("[Rest] Select error: ", err.Error())
	}
	return model.ValidateDuplicatedTag(name, tags)
}
