package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"repotags/handler"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetReposInternalServerError(t *testing.T) {
	req, err := http.NewRequest("GET", "/repositories", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler.GetRepos)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status == http.StatusInternalServerError {
		t.Log("DB not configured")
	}
}

func TestGetRepos(t *testing.T) {
	req, err := http.NewRequest("GET", "/repositories", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	rr.WriteHeader(http.StatusOK)
	handler := http.HandlerFunc(handler.GetRepos)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetReposByTagNullBody(t *testing.T) {
	req, err := http.NewRequest("GET", "/repositories/tag", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler.GetReposByTag)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status == http.StatusBadRequest {
		t.Log("Null body")
	}
}

func TestGetReposByTag(t *testing.T) {
	body := []byte(`{
				"name": "test"
				}`)
	req, err := http.NewRequest("GET", "/repositories/tag", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	rr.WriteHeader(http.StatusOK)
	handler := http.HandlerFunc(handler.GetReposByTag)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetReposByTagEmptyBody(t *testing.T) {
	body := []byte(`{
				"name": ""
				}`)
	req, err := http.NewRequest("GET", "/repositories/tag", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	rr.WriteHeader(http.StatusOK)
	handler := http.HandlerFunc(handler.GetReposByTag)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status == http.StatusBadRequest {
		t.Log("Empty body")
	}
}

func TestAddTagToRepo(t *testing.T) {
	body := []byte(`{
		"name": "test"
		}`)
	req, err := http.NewRequest("POST", "/repository/{id}/tag", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	urlvar := make(map[string]string)
	urlvar["id"] = "1234"
	mux.SetURLVars(req, urlvar)
	rr := httptest.NewRecorder()
	rr.WriteHeader(http.StatusOK)
	handler := http.HandlerFunc(handler.AddTagToRepo)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestAddTagToRepoNullId(t *testing.T) {
	body := []byte(`{
		"name": "test"
		}`)
	req, err := http.NewRequest("POST", "/repository/{id}/tag", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler.AddTagToRepo)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status == http.StatusBadRequest {
		t.Log("No repo ID!")
	}
}

func TestAddTagToRepoNullBody(t *testing.T) {
	req, err := http.NewRequest("POST", "/repository/{id}/tag", nil)
	if err != nil {
		t.Fatal(err)
	}
	urlvar := make(map[string]string)
	urlvar["id"] = "1234"
	mux.SetURLVars(req, urlvar)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler.AddTagToRepo)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status == http.StatusBadRequest {
		t.Log("No repo ID!")
	}
}
