package utils

import (
	"RJD02/job-portal/models"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func HandleResponse(w http.ResponseWriter, response models.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.ResponseCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON response", response.ResponseCode)
		return
	}
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err

}
