package utils

import (
	"RJD02/job-portal/models"
	"encoding/json"
	"net/http"
)

func handleResponse(w http.ResponseWriter, response models.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.ResponseCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON response", response.ResponseCode)
		return
	}
}
