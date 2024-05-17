package api

import (
	"log"
	"net/http"

	"github.com/DaffaJatmiko/go-rest-project-manager/repository"
	"github.com/DaffaJatmiko/go-rest-project-manager/service"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr  string
	store repository.Store
}

func NewAPIServer(addr string, store repository.Store) *APIServer {
	return &APIServer{addr: addr, store: store}
}

func (s *APIServer) Serve() {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	// registering services
	userService := service.NewUserService(s.store)
	userService.RegisterRoutes(subRouter)

	taskService := service.NewTaskService(s.store)
	taskService.RegisterRoutes(subRouter)

	projectService := service.NewProjectService(s.store)
	projectService.RegisterRoutes(subRouter)

	log.Println("Starting the API server at", s.addr)
	log.Fatal(http.ListenAndServe(s.addr, subRouter))
}