package service

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/DaffaJatmiko/go-rest-project-manager/config"
	"github.com/DaffaJatmiko/go-rest-project-manager/middleware"
	"github.com/DaffaJatmiko/go-rest-project-manager/model"
	"github.com/DaffaJatmiko/go-rest-project-manager/repository"
	"github.com/DaffaJatmiko/go-rest-project-manager/utils"
	"github.com/gorilla/mux"
)

var errEmailRequired = errors.New("email is required")
var errFirstNameRequired = errors.New("first name is required")
var errLastNameRequired = errors.New("last name is required")
var errPasswordRequired = errors.New("password is required")

type UserService struct {
	store repository.Store
}

func NewUserService(s repository.Store) *UserService {
	return &UserService{store: s}
}

func (s *UserService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users/register", s.handleUserRegister).Methods("POST")
	r.HandleFunc("/users/login", s.handleUserLogin).Methods("POST")
}

func (s *UserService) handleUserRegister(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var payload *model.User
	err = json.Unmarshal(body, &payload)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	if err := validateUserPayload(payload); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	hashedPW, err := middleware.HashPassword(payload.Password)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	payload.Password = hashedPW

	log.Printf("Registering user: %s, Hashed Password: %s", payload.Email, hashedPW)

	user, err := s.store.CreateUser(payload)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, model.ErrorResponse{Error: "error creating user"})
		return
	}

	token, err := createAndSetAuthCookie(user.ID, w)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, model.ErrorResponse{Error: "error creating token session"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, token)
}

func (s *UserService) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var payload *model.LoginRequest
	err = json.Unmarshal(body, &payload)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	if payload.Email == "" {
		utils.WriteJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: "missing email"})
		return
	}

	user, err := s.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("Login attempt for user: %s, Stored Hashed Password: %s", user.Email, user.Password)

	if !middleware.CheckPasswordHash(payload.Password, user.Password) {
		log.Printf("Password mismatch: Provided Password: %s, Hashed Password: %s", payload.Password, user.Password)
		utils.WriteJSON(w, http.StatusUnauthorized, model.ErrorResponse{Error: "wrong password"})
		return
	}

	token, err := createAndSetAuthCookie(user.ID, w)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, model.ErrorResponse{Error: "error creating token session"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, token)
}

func validateUserPayload(user *model.User) error {
	if user.Email == "" {
		return errEmailRequired
	}

	if user.FirstName == "" {
		return errFirstNameRequired
	}

	if user.LastName == "" {
		return errLastNameRequired
	}

	if user.Password == "" {
		return errPasswordRequired
	}

	return nil
}

func createAndSetAuthCookie(id int64, w http.ResponseWriter) (string, error) {
	secretKey := []byte(config.Envs.JWTSecret)
	token, err := middleware.CreateJWT(secretKey, id)
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	return token, nil
}
