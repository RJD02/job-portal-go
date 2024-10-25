package controllers

import (
	"RJD02/job-portal/config"
	"RJD02/job-portal/db"
	"RJD02/job-portal/models"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	SECRET_KEY := config.AppConfig.JWT_SECRET_KEY
	var req_user models.User
	var response models.Response
	var user *db.UserModel

	ctx := context.Background()

	err := json.NewDecoder(r.Body).Decode(&req_user)
	if err != nil {
		response.ResponseCode = http.StatusBadRequest
		response.Message = "Error in decoding"
		response.Error = err.Error()
		handleResponse(w, response)
		return
	}

	if req_user.Username == "" {
		user, err = config.AppConfig.Db.User.FindUnique(
			db.User.Email.Equals(req_user.Email),
		).Exec(ctx)
		log.Println("No username")
	} else if req_user.Email == "" {
		user, err = config.AppConfig.Db.User.FindUnique(
			db.User.Username.Equals(req_user.Username),
		).Exec(ctx)
		log.Println("No Email")
	} else {
		response.ResponseCode = http.StatusBadRequest
		response.Message = "Provide either username or email"
		handleResponse(w, response)
		return
	}

	if err != nil {
		response.ResponseCode = http.StatusBadRequest
		response.Message = "Both username and email didn't match"
		response.Error = err.Error()
		response.Data = nil
		handleResponse(w, response)
		return
	}

	if user != nil {
		userJson, err := json.Marshal(user)
		if err != nil {

		} else {
			log.Println(string(userJson))
		}
	}

	// hash req_user's password and check with db_user's password
	if !checkPasswordHash(req_user.Password, user.Password) {
		response.Message = "Wrong password provided"
		response.ResponseCode = http.StatusBadRequest
		handleResponse(w, response)
	}

	// jwt logic
	expiry := time.Now().Add(time.Hour * 24 * 2)

	claims := jwt.MapClaims{
		"username": user.Username,
		"exp":      expiry.Unix(),
		"email":    user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		response.Error = err.Error()
		response.ResponseCode = http.StatusInternalServerError
		response.Message = "Error while siging the jwt token"
		handleResponse(w, response)
		return
	}

	updatedUser, err := config.AppConfig.Db.User.FindUnique(
		db.User.ID.Equals(user.ID),
	).Update(db.User.Token.Set(signedToken), db.User.Expiry.Set(expiry)).Exec(ctx)

	if err != nil {
		response.Error = err.Error()
		response.ResponseCode = http.StatusInternalServerError
		response.Message = "Error while updating jwt token and expiry date for this user"
		handleResponse(w, response)
		return
	}

	response.Data = updatedUser
	response.Message = "Successfully fetched the user"
	response.ResponseCode = http.StatusOK
	handleResponse(w, response)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	SECRET_KEY := config.AppConfig.JWT_SECRET_KEY
	var req_user models.User
	ctx := context.Background()
	response := models.Response{}
	err := json.NewDecoder(r.Body).Decode(&req_user)
	if err != nil {
		response.Message = "Something went wrong while decoding"
		response.Error = err.Error()
		response.ResponseCode = http.StatusBadRequest
		log.Println("Error while decoding request", err)
		handleResponse(w, response)
		return
	}

	log.Println("Request body extracted: ", req_user.Username, req_user.Email, req_user.Password)

	if req_user.Password == "" || req_user.Email == "" || req_user.Username == "" {
		response.ResponseCode = http.StatusBadRequest
		response.Error = "Not all data is present"
		response.Message = "Data is not correct"
		handleResponse(w, response)
		return
	}

	req_user.Password, err = hashPassword(req_user.Password)
	if err != nil {
		response.ResponseCode = http.StatusInternalServerError
		response.Message = "Error while password hashing"
		response.Error = err.Error()
		handleResponse(w, response)
		return
	}

	expiry := time.Now().Add(time.Hour * 24 * 2)

	claims := jwt.MapClaims{
		"username": req_user.Username,
		"exp":      expiry.Unix(),
		"email":    req_user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		response.Error = err.Error()
		response.ResponseCode = http.StatusInternalServerError
		response.Message = "Error while signing the jwt token"
		handleResponse(w, response)
		return
	}

	user, err := config.AppConfig.Db.User.CreateOne(
		db.User.Username.Set(req_user.Username),
		db.User.Password.Set(req_user.Password),
		db.User.Email.Set(req_user.Email),
		db.User.Token.Set(signedToken),
		db.User.Expiry.Set(expiry),
	).Exec(ctx)

	if err != nil {
		response.ResponseCode = http.StatusInternalServerError
		response.Message = "Error while inserting new user to db"
		response.Error = err.Error()
		handleResponse(w, response)
		return
	}

	response.ResponseCode = http.StatusOK
	response.Message = "User Added"
	response.Data = user
	handleResponse(w, response)

}
