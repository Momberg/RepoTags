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

//GetReposTagRecommendation get repositories tags recommendation
func GetReposTagRecommendation(w http.ResponseWriter, r *http.Request) {
	db, err := repository.GetDBConnection()
	if err != nil {
		log.Println("[Rest] Connection error: ", err.Error())
		return
	}
	sql := "select id, name, description, url, language from repository"
	repos := []model.Repository{}
	err = db.Select(&repos, sql)
	if err != nil {
		log.Println("[Rest] Select error: ", err.Error())
	}
	getTagRecommendation(repos)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repos)
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
	tagMap := make(map[string]int64)
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
		tagMap[tag.Name] = repo.ID
		repos = append(repos, repo)
	}
	removeDBDuplicated(repos)
	for index := range repos {
		tagS := model.Tags{}
		for tag, id := range tagMap {
			if id == repos[index].ID {
				tagS.Name = tag
				repos[index].Tags = append(repos[index].Tags, tagS)
			}
		}
	}
	getTagRecommendation(repos)
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
	if validateTags(tag.Name, params["id"]) {
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
func validateTags(name string, id string) bool {
	db, err := repository.GetDBConnection()
	if err != nil {
		log.Println("[Rest] Connection error: ", err.Error())
	}
	sql := "select id, name from tags"
	tags := []model.Tags{}
	err = db.Select(&tags, sql)
	if err != nil {
		log.Println("[Rest] Tags select error: ", err.Error())
	}
	return model.ValidateDuplicatedTag(name, tags, id)
}

//removeDBDuplicated remove duplicated data
func removeDBDuplicated(repos []model.Repository) {
	index := 0
	for range repos {
		if len(repos)-1 > index {
			if repos[index].ID == repos[index+1].ID {
				repos[index] = repos[len(repos)-1]
				repos[len(repos)-1] = model.Repository{}
				repos = repos[:len(repos)-1]
				index--
			}
		}
		index++
	}
}

//getTagRecommendation recommend tags
func getTagRecommendation(repos []model.Repository) {
	for repo := range repos {
		if repos[repo].Language != "" {
			repos[repo].Tagrecommendation = repos[repo].Language
		} else {
			repos[repo].Tagrecommendation = repos[repo].Name
		}
	}
}
