package controllers

import (
	"RJD02/job-portal/config"
	"RJD02/job-portal/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func AuthHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You've hit api route"))
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err

}

func Login(w http.ResponseWriter, r *http.Request) {
	var req_user models.User
	var db_user models.User

	err := json.NewDecoder(r.Body).Decode(&req_user)
	if err != nil {
		http.Error(w, "Error in decoding", http.StatusBadRequest)
		return
	}

	err = config.AppConfig.Db.QueryRow("SELECT username, password FROM job_portal.users WHERE username = $1", req_user.Username).Scan(&db_user.Username, &db_user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error in getting the user", http.StatusInternalServerError)
			log.Println("Error in getting the user", err)
		}
		return
	}

	// hash req_user's password and check with db_user's password
	if checkPasswordHash(req_user.Password, db_user.Password) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(req_user)
	} else {
		log.Println("Wrong password provided", err)
		http.Error(w, "Wrong password", http.StatusBadRequest)
		return
	}
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var req_user models.User

	err := json.NewDecoder(r.Body).Decode(&req_user)
	if err != nil {
		log.Println("Error while decoding request", err)
		http.Error(w, "Something went wrong while decoding", http.StatusBadRequest)
		return
	}

	req_user.Password, err = hashPassword(req_user.Password)
	if err != nil {
		log.Println("Error while password hashing", err)
		http.Error(w, "Something went wrong while hashing password", http.StatusInternalServerError)
		return
	}

	result, err := config.AppConfig.Db.Exec(`
        insert into job_portal.users (
            username,
            password
        ) values ($1, $2)
    `, req_user.Username, req_user.Password)

	if err != nil {
		log.Println("Error while inserting new user to db", err)
		http.Error(w, "Db error", http.StatusInternalServerError)
		return
	}

	rows_affected, err := result.RowsAffected()

	last_inserted_id, err := result.LastInsertId()

	log.Println(last_inserted_id, "Added to db, ", rows_affected, "rows affected")

	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, "User created")
}
