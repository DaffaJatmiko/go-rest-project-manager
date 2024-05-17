package repository

import (
	"github.com/DaffaJatmiko/go-rest-project-manager/model"
)

// MockStore implements the Store interface for testing purposes
type MockStore struct{}

func (s *MockStore) CreateProject(p *model.Project) (*model.Project, error) {
	return p, nil
}

func (s *MockStore) GetProject(id string) (*model.Project, error) {
	return &model.Project{Name: "Super cool project"}, nil
}

func (s *MockStore) DeleteProject(id string) error {
	return nil
}

func (s *MockStore) UpdateProject(project *model.Project) (*model.Project, error) {
	return project, nil
}

func (s *MockStore) CreateUser(u *model.User) (*model.User, error) {
	return u, nil
}

func (s *MockStore) GetUserByID(id string) (*model.User, error) {
	return &model.User{}, nil
}

func (s *MockStore) GetUserByEmail(email string) (*model.User, error) {
	return &model.User{}, nil
}

func (s *MockStore) CreateTask(t *model.Task) (*model.Task, error) {
	return t, nil
}

func (s *MockStore) GetTask(id string) (*model.Task, error) {
	return &model.Task{}, nil
}

func (s *MockStore) DeleteTask(id string) error {
	return nil
}

func (s *MockStore) UpdateTask(task *model.Task) (*model.Task, error) {
	return task, nil
}
