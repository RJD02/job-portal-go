package utils

import (
	"RJD02/job-portal/config"
	"RJD02/job-portal/db"
	"RJD02/job-portal/models"
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func updateUser(user db.UserModel, signedToken string, expiry time.Time) (*db.UserModel, error) {
	ctx := context.Background()
	updatedUser, err := config.AppConfig.Db.User.FindUnique(
		db.User.ID.Equals(user.ID),
	).Update(db.User.Token.Set(signedToken), db.User.Expiry.Set(expiry)).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func CreateTokenAndUpdateUser(user db.UserModel) (string, *db.UserModel, error) {
	SECRET_KEY := config.AppConfig.JWT_SECRET_KEY

	expiry := time.Now().Add(time.Hour * 24 * 2)

	claims := jwt.MapClaims{
		"username": user.Username,
		"exp":      expiry.Unix(),
		"email":    user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", nil, err
	}

	updatedUser, err := updateUser(user, signedToken, expiry)

	if err != nil {
		return "", nil, err
	}

	return signedToken, updatedUser, nil
}

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
