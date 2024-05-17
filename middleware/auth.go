package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/DaffaJatmiko/go-rest-project-manager/config"
	"github.com/DaffaJatmiko/go-rest-project-manager/model"
	"github.com/DaffaJatmiko/go-rest-project-manager/repository"
	"github.com/DaffaJatmiko/go-rest-project-manager/utils"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func AuthHandler(handlerFunc http.HandlerFunc, store repository.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// get the token from the request header (auth)
		tokenStr := GetToken(r)
		// validate the token
		token, err := ValidateJWT(tokenStr)
		if err != nil {
			log.Println(err)
			permmissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("token is not valid")
			permmissionDenied(w)
			return
		}
		// get the userId from the token with claims
		claims := token.Claims.(jwt.MapClaims)
		userID := claims["userID"].(string)

		_, err = store.GetUserByID(userID)
		if err != nil {
			log.Println("failed to get user")
			permmissionDenied(w)
			return
		}
		// call the handler func and continue to the next endpoint
		handlerFunc(w, r)

	}
}

func permmissionDenied(w http.ResponseWriter) {
	utils.WriteJSON(w, http.StatusUnauthorized, model.ErrorResponse{
		Error: fmt.Errorf("permission denied").Error(),
	})
}

func GetToken(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}

func ValidateJWT(token string) (*jwt.Token, error) {
	secretKey := config.Envs.JWTSecret

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateJWT(secretKey []byte, userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": strconv.Itoa(int(userID)),
		"expiresAt": time.Now().Add(time.Minute * 1).Unix(), 
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
