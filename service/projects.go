package service

import (
	"database/sql"
	"encoding/json"
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

type ProjectService struct {
	store repository.Store
}

func NewProjectService(s repository.Store) *ProjectService {
	return &ProjectService{store: s}
}

func (s *ProjectService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/projects", middleware.AuthHandler(s.handleCreateProject, s.store)).Methods("POST")
	r.HandleFunc("/projects/{id}", middleware.AuthHandler(s.handleGetProject, s.store)).Methods("GET")
	r.HandleFunc("/projects/{id}", middleware.AuthHandler(s.handleDeleteProject, s.store)).Methods("DELETE")
	r.HandleFunc("/projects/{id}", middleware.AuthHandler(s.handleUpdateProject, s.store)).Methods("PUT")
}


func (s *ProjectService) handleCreateProject(w http.ResponseWriter, r *http.Request){
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: "invalid request payload"})
		return
	}

	defer r.Body.Close()

	var project *model.Project
	if err := json.Unmarshal(body, &project); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: "invalid request payload"})
		return
	}

	if project.Name == "" {
		utils.WriteJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: "missing name"})
		return
	}

	p, err := s.store.CreateProject(project)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, p)

}

func (s *ProjectService) handleGetProject(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]

	project, err := s.store.GetProject(id)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
	}

	utils.WriteJSON(w, http.StatusOK, project)
}

func (s *ProjectService) handleDeleteProject(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]

	log.Printf("handleDeleteProject: attempting to delete project id %s", id)
	
	err := s.store.DeleteProject(id)
	if err != nil {
			if err == sql.ErrNoRows {
					log.Printf("handleDeleteProject: project id %s not found", id)
					utils.WriteJSON(w, http.StatusNotFound, model.ErrorResponse{Error: "Project not found"})
					return
			}
			log.Printf("handleDeleteProject: error deleting project id %s: %v", id, err)
			utils.WriteJSON(w, http.StatusInternalServerError, model.ErrorResponse{Error: "Error deleting project"})
			return
	}

	log.Printf("handleDeleteProject: project id %s deleted successfully", id)
	utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("Deleted project %s", id))
}

func (s *ProjectService) handleUpdateProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: "invalid request payload"})
		return
	}

	defer r.Body.Close()

	var project *model.Project
	if err := json.Unmarshal(body, &project); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: "invalid request payload"})
		return
	}

	project.ID, _ = strconv.ParseInt(id, 10, 64)
	if project.Name == "" {
		utils.WriteJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: "missing name"})
		return
	}

	updatedProject, err := s.store.UpdateProject(project)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, model.ErrorResponse{Error: "failed to update project"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, updatedProject)

}
