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

func TestCreateProject(t *testing.T) {
	memoryStore := &repository.MockStore{}
	service := NewProjectService(memoryStore)

	t.Run("should return an error when name is empty", func(t *testing.T) {
			payload := &model.Project{
					Name: "",
			}

			b, err := json.Marshal(payload)
			if err != nil {
					t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/projects", bytes.NewBuffer(b))
			if err != nil {
					t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			router := mux.NewRouter()

			router.HandleFunc("/projects", service.handleCreateProject)

			router.ServeHTTP(rr, req)

			if rr.Code != http.StatusBadRequest {
					t.Error("invalid status code, it should fail")
			}
	})

	t.Run("should create a project", func(t *testing.T) {
			payload := &model.Project{
					Name: "New Project",
			}

			b, err := json.Marshal(payload)
			if err != nil {
					t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/projects", bytes.NewBuffer(b))
			if err != nil {
					t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()

			router.HandleFunc("/projects", service.handleCreateProject)

			router.ServeHTTP(rr, req)

			if rr.Code != http.StatusCreated {
					t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
			}
	})
}

func TestGetProject(t *testing.T) {
	memoryStore := &repository.MockStore{}
	service := NewProjectService(memoryStore)

	t.Run("should return the project by id", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/projects/42", nil)
			if err != nil {
					t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()

			router.HandleFunc("/projects/{id}", service.handleGetProject)

			router.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
					t.Error("invalid status code", rr.Code)
			}
	})
}

func TestDeleteProject(t *testing.T) {
	memoryStore := &repository.MockStore{}
	service := NewProjectService(memoryStore)

	t.Run("should delete the project by id", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, "/projects/42", nil)
			if err != nil {
					t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()

			router.HandleFunc("/projects/{id}", service.handleDeleteProject)

			router.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
					t.Error("invalid status code", rr.Code)
			}
	})
}

func TestUpdateProject(t *testing.T) {
	memoryStore := &repository.MockStore{}
	service := NewProjectService(memoryStore)

	t.Run("should update the project by id", func(t *testing.T) {
			payload := &model.Project{
					Name: "Updated Project",
			}

			b, err := json.Marshal(payload)
			if err != nil {
					t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPut, "/projects/42", bytes.NewBuffer(b))
			if err != nil {
					t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()

			router.HandleFunc("/projects/{id}", service.handleUpdateProject)

			router.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
					t.Error("invalid status code", rr.Code)
			}
	})
}