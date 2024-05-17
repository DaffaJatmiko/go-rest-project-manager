package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DaffaJatmiko/go-rest-project-manager/model"
	"github.com/DaffaJatmiko/go-rest-project-manager/repository"
	"github.com/gorilla/mux"
)

func TestCreateTask(t *testing.T) {
	memoryStore := &repository.MockStore{}
	service := NewTaskService(memoryStore)

	t.Run("should return an error when name is empty", func(t *testing.T) {
		payload := &model.Task{
			Name: "",
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/tasks", service.handleCreateTask).Methods("POST")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should create a task", func(t *testing.T) {
		payload := &model.Task{
			Name:        "Creating a REST API in go",
			ProjectID:   1,
			AssignedToID: 42,
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/tasks", service.handleCreateTask).Methods("POST")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

func TestGetTask(t *testing.T) {
	memoryStore := &repository.MockStore{}
	service := NewTaskService(memoryStore)

	t.Run("should return the task by id", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/tasks/42", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/tasks/{id}", service.handleGetTask).Methods("GET")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

func TestDeleteTask(t *testing.T) {
	memoryStore := &repository.MockStore{}
	service := NewTaskService(memoryStore)

	t.Run("should delete the task by id", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, "/tasks/42", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/tasks/{id}", service.handleDeleteTask).Methods("DELETE")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

func TestUpdateTask(t *testing.T) {
	memoryStore := &repository.MockStore{}
	service := NewTaskService(memoryStore)

	t.Run("should update the task", func(t *testing.T) {
		payload := &model.Task{
			ID:          42,
			Name:        "Updated task",
			ProjectID:   1,
			AssignedToID: 3,
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPut, "/tasks/42", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/tasks/{id}", service.handleUpdateTask).Methods("PUT")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}