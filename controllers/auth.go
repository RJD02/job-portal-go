package controllers

import (
	"RJD02/job-portal/config"
	"RJD02/job-portal/db"
	"RJD02/job-portal/mail"
	"RJD02/job-portal/models"
	"RJD02/job-portal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"context"
	"sync"

	"github.com/golang-jwt/jwt/v5"
)

func AuthHome(w http.ResponseWriter, r *http.Request) {
	var response models.Response
	response.ResponseCode = http.StatusOK
	response.Message = "You've hit api route"
	utils.HandleResponse(w, response)
}

func MagicLogin(w http.ResponseWriter, r *http.Request) {
	queryToken := r.URL.Query().Get("token")
	queryEmail := r.URL.Query().Get("email")
	var response models.Response
	// var tokenEmail string

	token, err := jwt.Parse(queryToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.AppConfig.JWT_SECRET_KEY), nil
	})

	if err != nil {
		response.ResponseCode = http.StatusBadRequest
		response.Error = err.Error()
		response.Message = "Error parsing token"
		utils.HandleResponse(w, response)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		for key, val := range claims {
			fmt.Printf("%s: %v\n", key, val)
		}
	} else {
		response.ResponseCode = http.StatusBadRequest
		response.Message = "Invalid token"
		utils.HandleResponse(w, response)
		return
	}

	// basic validations
	user, err := config.AppConfig.Db.User.FindUnique(
		db.User.Email.Equals(queryEmail),
	).Exec(context.Background())

	if err != nil {
		response.ResponseCode = http.StatusBadRequest
		response.Message = "User not found"
		response.Error = err.Error()
		utils.HandleResponse(w, response)
		return
	}

	response.ResponseCode = http.StatusOK
	response.Message = "Login Successful"
	response.Data = user
	utils.HandleResponse(w, response)

}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req_user models.User
	var response models.Response
	var user *db.UserModel
	var wg sync.WaitGroup
	// ch := make(chan string)
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&req_user)
	ctx := context.Background()

	if err != nil {
		response.Error = err.Error()
		response.Message = "Something went wrong while decoding body"
		response.ResponseCode = http.StatusBadRequest
		utils.HandleResponse(w, response)
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
		utils.HandleResponse(w, response)
		return
	}

	if err != nil {
		response.ResponseCode = http.StatusBadRequest
		response.Message = "Both username and email didn't match"
		response.Error = err.Error()
		response.Data = nil
		utils.HandleResponse(w, response)
		return
	}

	// create token and attach to the user
	signedToken, updatedUser, err := utils.CreateTokenAndUpdateUser(*user)
	if err != nil {
		response.ResponseCode = http.StatusInternalServerError
		response.Message = "Error while creating a token"
		response.Error = err.Error()
		utils.HandleResponse(w, response)
		return
	}
	baseURL := r.Host
	magicLink := fmt.Sprintf("http://%s/auth/magic-login?email=%s&token=%s", baseURL, updatedUser.Email, signedToken)
	emailBody := mail.GenerateMagicLinkEmail(updatedUser.Username, magicLink)
	log.Println("Email body: ", emailBody)
	subject := "Here's your magic link to login"

	wg.Add(1)
	go mail.SendMail(updatedUser.Email, emailBody, subject, &wg)

	wg.Wait()
	response.ResponseCode = http.StatusOK
	response.Message = "Successfully sent to registered email"
	utils.HandleResponse(w, response)

}

func Login(w http.ResponseWriter, r *http.Request) {
	var req_user models.User
	var response models.Response
	var user *db.UserModel

	ctx := context.Background()
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&req_user)
	if err != nil {
		response.ResponseCode = http.StatusBadRequest
		response.Message = "Error in decoding"
		response.Error = err.Error()
		utils.HandleResponse(w, response)
		return
	}

	if req_user.Email != "" {
		user, err = config.AppConfig.Db.User.FindUnique(
			db.User.Email.Equals(req_user.Email),
		).Exec(ctx)
		log.Println("No username")
	} else if req_user.Username != "" {
		user, err = config.AppConfig.Db.User.FindUnique(
			db.User.Username.Equals(req_user.Username),
		).Exec(ctx)
		log.Println("No Email")
	} else {
		response.ResponseCode = http.StatusBadRequest
		response.Message = "Provide either username or email"
		utils.HandleResponse(w, response)
		return
	}

	// hash req_user's password and check with db_user's password
	if !utils.CheckPasswordHash(req_user.Password, user.Password) {
		response.Message = "Wrong password provided"
		response.ResponseCode = http.StatusBadRequest
		utils.HandleResponse(w, response)
	}

	_, updatedUser, err := utils.CreateTokenAndUpdateUser(*user)

	if err != nil {
		response.Error = err.Error()
		response.ResponseCode = http.StatusInternalServerError
		response.Message = "Error while updating jwt token and expiry date for this user"
		utils.HandleResponse(w, response)
		return
	}

	response.Data = updatedUser
	response.Message = "Successfully fetched the user"
	response.ResponseCode = http.StatusOK
	utils.HandleResponse(w, response)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var req_user models.User
	ctx := context.Background()
	response := models.Response{}
	err := json.NewDecoder(r.Body).Decode(&req_user)
	if err != nil {
		response.Message = "Something went wrong while decoding"
		response.Error = err.Error()
		response.ResponseCode = http.StatusBadRequest
		log.Println("Error while decoding request", err)
		utils.HandleResponse(w, response)
		return
	}

	if req_user.Password == "" || req_user.Email == "" || req_user.Username == "" {
		response.ResponseCode = http.StatusBadRequest
		response.Error = "Not all data is present"
		response.Message = "Data is not correct"
		utils.HandleResponse(w, response)
		return
	}

	req_user.Password, err = utils.HashPassword(req_user.Password)
	if err != nil {
		response.ResponseCode = http.StatusInternalServerError
		response.Message = "Error while password hashing"
		response.Error = err.Error()
		utils.HandleResponse(w, response)
		return
	}

	user, err := config.AppConfig.Db.User.CreateOne(
		db.User.Username.Set(req_user.Username),
		db.User.Password.Set(req_user.Password),
		db.User.Email.Set(req_user.Email),
	).Exec(ctx)

	if err != nil {
		if _, ok := db.IsErrUniqueConstraint(err); ok {
			response.ResponseCode = http.StatusBadRequest
			response.Message = "Username/Email already exists"
			response.Error = err.Error()
			utils.HandleResponse(w, response)
			return
		}
		response.ResponseCode = http.StatusInternalServerError
		response.Message = "Error while inserting new user to db"
		response.Error = err.Error()
		utils.HandleResponse(w, response)
		return
	}

	_, updatedUser, err := utils.CreateTokenAndUpdateUser(*user)

	if err != nil {
		response.ResponseCode = http.StatusInternalServerError
		response.Message = "Error while create and update user"
		response.Error = err.Error()
		utils.HandleResponse(w, response)
		return
	}

	response.ResponseCode = http.StatusOK
	response.Message = "User Added"
	response.Data = updatedUser
	utils.HandleResponse(w, response)

}
