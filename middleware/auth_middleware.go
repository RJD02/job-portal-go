package middleware

import (
	"RJD02/job-portal/config"
	"RJD02/job-portal/models"
	"RJD02/job-portal/utils"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var response models.Response
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			response.ResponseCode = http.StatusUnauthorized
			response.Message = "Missing token"
			utils.HandleResponse(w, response)
			return
		}

		splitToken := strings.Split(tokenString, "Bearer ")

		if len(splitToken) != 2 {
			response.ResponseCode = http.StatusUnauthorized
			response.Message = "Invalid token format"
			utils.HandleResponse(w, response)
			return
		}

		tokenString = splitToken[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.AppConfig.JWT_SECRET_KEY), nil
		})

		if err != nil || !token.Valid {
			response.ResponseCode = http.StatusUnauthorized
			response.Message = "Invalid token"
			utils.HandleResponse(w, response)
			return
		}

		// Add token claims to the context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if exp, ok := claims["exp"].(float64); ok && time.Unix(int64(exp), 0).Before(time.Now()) {
				response.ResponseCode = http.StatusUnauthorized
				response.Message = "Token expired"
				utils.HandleResponse(w, response)
				return
			}

			ctx := context.WithValue(r.Context(), "user", claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			response.ResponseCode = http.StatusUnauthorized
			response.Message = "Invalid token claims"
			utils.HandleResponse(w, response)
			return
		}
	})
}
