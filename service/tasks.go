package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/DaffaJatmiko/go-rest-project-manager/middleware"
	"github.com/DaffaJatmiko/go-rest-project-manager/model"
	"github.com/DaffaJatmiko/go-rest-project-manager/repository"
	"github.com/DaffaJatmiko/go-rest-project-manager/utils"
	"github.com/gorilla/mux"
)

var errNameRequired = errors.New("name is required")
var errProjectIDRequired = errors.New("project id is required")
var errUserIDRequired = errors.New("user id is required")

type TaskService struct {
	store repository.Store
}

func NewTaskService(s repository.Store) *TaskService {
	return &TaskService{store: s}
}

func (s *TaskService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/tasks", middleware.AuthHandler(s.handleCreateTask, s.store)).Methods("POST")
	r.HandleFunc("/tasks/{id}", middleware.AuthHandler(s.handleGetTask, s.store)).Methods("GET")
	r.HandleFunc("/tasks/{id}", middleware.AuthHandler(s.handleDeleteTask, s.store)).Methods("DELETE")
	r.HandleFunc("/tasks/{id}", middleware.AuthHandler(s.handleUpdateTask, s.store)).Methods("PUT")
}

func (s *TaskService) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, model.ErrorResponse{Error: "Invalid request payload"})
			return
	}
	defer r.Body.Close()

	var task model.Task
	if err = json.Unmarshal(body, &task); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request payload"})
			return
	}

	if err := validateTaskPayload(&task); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
	}

	log.Printf("Creating task: %+v", task)

	t, err := s.store.CreateTask(&task)
	if err != nil {
			log.Printf("Error creating task: %v", err)
			utils.WriteJSON(w, http.StatusInternalServerError, model.ErrorResponse{Error: "Error creating task"})
			return
	}

	utils.WriteJSON(w, http.StatusCreated, t)
}


func (s *TaskService) handleGetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]


	task, err := s.store.GetTask(id)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, task)
}

func (s *TaskService) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := s.store.DeleteTask(id); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, model.ErrorResponse{Error: "error deleting task"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("Deleted task %s", id))

}

func (s *TaskService) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, model.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	defer r.Body.Close()

	var task *model.Task
	if err := json.Unmarshal(body, &task); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: "invalid request payload"})
		return
	}

	task.ID, _ = strconv.ParseInt(id, 10, 64)
	if err := validateTaskPayload(task); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	updatedTask, err := s.store.UpdateTask(task)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, model.ErrorResponse{Error: "Error updating task"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, updatedTask)
}

func validateTaskPayload(task *model.Task) error {
	if task.Name == "" {
		return errNameRequired
	}

	if task.ProjectID == 0 {
		return errProjectIDRequired
	}

	if task.AssignedToID == 0 {
		return errUserIDRequired
	}
	return nil
}