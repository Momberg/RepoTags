package handler

import (
	"fmt"
	"net/http"
)

//CreateRestHandler responsable of create the REST call
func CreateRestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Creating handler")
}
