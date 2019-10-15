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

//GetReposByTag get repositories
func GetReposByTag(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if len(params["id"]) < 1 {
		log.Println("[Rest] Missing repository tag name.")
		return
	}
	db, err := repository.GetDBConnection()
	if err != nil {
		log.Println("[Rest] Connection error: ", err.Error())
		return
	}
	sql := "select r.id, r.name, r.description, r.url, r.language, t.name as tags from repository r join tags t where r.tag_id = t.id and t.name like ? GROUP BY r.id, t.name;"
	repos := []model.Repository{}
	repo := model.Repository{}
	rows, _ := db.Query(sql, params["id"]+"%")
	for rows.Next() {
		tag := model.Tags{}
		err = rows.Scan(&repo.ID,
			&repo.Name,
			&repo.Description,
			&repo.URL,
			&repo.Language,
			&tag.Name)
		if err != nil {
			log.Println("[Rest] Scan error: ", err.Error())
		}
		repo.Tags = append(repo.Tags, tag)
		repos = append(repos, repo)
	}
	removeDBDuplicated(repos)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repos)
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

//removeDBDuplicated remove duplicated data
func removeDBDuplicated(repos []model.Repository) {
	index := 0
	for range repos {
		if len(repos)-1 > index {
			if repos[index].ID == repos[index+1].ID {
				/*if len(repos[index].Tags) < len(repos[index+1].Tags) {
					repos[index] = repos[len(repos)-1]
					repos[len(repos)-1] = model.Repository{}
					repos = repos[:len(repos)-1]
					index--
					continue
				}*/
				repos[index] = repos[len(repos)-1]
				repos[len(repos)-1] = model.Repository{}
				repos = repos[:len(repos)-1]
				index--
			}
		}
		index++
	}
}
