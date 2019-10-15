package handler_test

import (
	"net/http"
	"net/http/httptest"
	"repotags/handler"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetRepos(t *testing.T) {
	req, err := http.NewRequest("GET", "/repositories", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler.GetRepos)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetReposByTag(t *testing.T) {
	req, err := http.NewRequest("GET", "/repositories/tag/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}
	urlvar := make(map[string]string)
	urlvar["id"] = "1234"
	mux.SetURLVars(req, urlvar)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler.GetReposByTag)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
