package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"repotags/model"
	"repotags/repository"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

//GetRepos get repositories
func GetRepos(w http.ResponseWriter, r *http.Request) {
	db, err := repository.GetDBConnection()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("[Rest] Connection error: ", err.Error())
		return
	}
	sql := "select id, name, description, url, language from repository"
	repos := []model.Repository{}
	err = db.Select(&repos, sql)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("[Rest] Select error: ", err.Error())
		return
	}
	getTagRecommendation(repos)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repos)
}

//GetReposByTag get repositories
func GetReposByTag(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tagBody := model.Tags{}
		err = json.Unmarshal(body, &tagBody)
		if len(tagBody.Name) < 1 || tagBody.Name == "" {
			http.Error(w, "Missing repository ID.", http.StatusBadRequest)
			log.Println("[Rest] Missing repository tag name.")
			return
		}
		db, err := repository.GetDBConnection()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("[Rest] Connection error: ", err.Error())
			return
		}
		sql := "select r.id, r.name, r.description, r.url, r.language, t.name as tags from repository r join tags t where r.tag_id = t.id and t.name like ? GROUP BY r.id, t.name order by r.id;"
		repos := []model.Repository{}
		repo := model.Repository{}
		rows, _ := db.Query(sql, tagBody.Name+"%")
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
			tagMap[tag.Name+"-"+strconv.FormatInt(repo.ID, 10)] = repo.ID
			repos = append(repos, repo)
		}
		repos = removeDuplicates(repos)
		for index := range repos {
			tagS := model.Tags{}
			for tag, id := range tagMap {
				if id == repos[index].ID {
					tagS.Name = tag[:strings.IndexByte(tag, '-')]
					repos[index].Tags = append(repos[index].Tags, tagS)
				}
			}
		}
		getTagRecommendation(repos)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(repos)
	}
	http.Error(w, "Body cannot be null.", http.StatusBadRequest)
	return
}

//AddTagToRepo add tag to repository
func AddTagToRepo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if len(params["id"]) < 1 {
		http.Error(w, "Missing repository ID.", http.StatusBadRequest)
		log.Println("[Rest] Missing repository ID.")
		return
	}
	if r.Body != nil {
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tag := model.Tags{}
		err = json.Unmarshal(body, &tag)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("[Rest] Body convertion error: ", err.Error())
			return
		}
		if validateTags(tag.Name, params["id"]) {
			http.Error(w, "This tag already exists for this repository.", http.StatusBadRequest)
			return
		}
		db, err := repository.GetDBConnection()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("[Rest] Connection error: ", err.Error())
			return
		}
		sql := "insert into tags (id, name) values (?, ?)"
		_, err = db.Exec(sql, params["id"], &tag.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("[Rest] Insert tags error: ", err.Error())
		}
		sql = "update repository set tag_id = (?) where id = " + params["id"]
		_, err = db.Exec(sql, params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("[Rest] Update repository error: ", err.Error())
		}
	}
	http.Error(w, "Body cannot be null.", http.StatusBadRequest)
	return
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

//removeDuplicates remove duplicated repesotories Ids
func removeDuplicates(repos []model.Repository) []model.Repository {
	encountered := map[int64]bool{}
	result := []model.Repository{}
	for v := range repos {
		if encountered[repos[v].ID] == true {
		} else {
			encountered[repos[v].ID] = true
			result = append(result, repos[v])
		}
	}
	return result
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
