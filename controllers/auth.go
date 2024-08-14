package controllers

import (
	"RJD02/job-portal/models"
	"encoding/json"
	"net/http"
)

func AuthHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You've hit api route"))
}

func Login(w http.ResponseWriter, r *http.Request) {
	var j models.User

	err := json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		http.Error(w, "Error in decoding", http.StatusBadRequest)
		return
	}

}
