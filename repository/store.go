package repository

import (
	"database/sql"
	"log"

	"github.com/DaffaJatmiko/go-rest-project-manager/model"
)

type Store interface {
	// Users
	CreateUser(u *model.User) (*model.User, error)
	GetUserByID(id string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	// Tasks
	CreateTask(t *model.Task) (*model.Task, error)
	GetTask(id string) (*model.Task, error)
	DeleteTask(id string) error
	UpdateTask(t *model.Task) (*model.Task, error)
	// Projects
	CreateProject(p *model.Project) (*model.Project, error)
	GetProject(id string) (*model.Project, error)
	DeleteProject(id string) error 
	UpdateProject(p *model.Project) (*model.Project, error)
}

type Storage struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateUser(u *model.User) (*model.User, error) {
	rows, err := s.db.Exec("INSERT INTO users (email, firstName, lastName, password) VALUES (?, ?, ?, ?)", u.Email, u.FirstName, u.LastName, u.Password)
	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	u.ID = id
	return u, nil
}

func (s *Storage) CreateTask(t *model.Task) (*model.Task, error) {
	query := "INSERT INTO tasks (name, status, projectId, assignedToID) VALUES (?, ?, ?, ?)"
	result, err := s.db.Exec(query, t.Name, t.Status, t.ProjectID, t.AssignedToID)
	if err != nil {
			log.Printf("CreateTask: error executing query: %v", err)
			return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
			log.Printf("CreateTask: error getting last insert ID: %v", err)
			return nil, err
	}

	t.ID = id
	return t, nil
}


func (s *Storage) CreateProject(p *model.Project) (*model.Project, error) {
	rows, err := s.db.Exec("INSERT INTO projects (name) values (?)", p.Name)
	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	p.ID = id 

	return p, nil
}

func (s *Storage) GetUserByID(id string) (*model.User, error) {
	var user model.User
	err := s.db.QueryRow("SELECT id, email, firstName, lastName, createdAt FROM users WHERE id = ?", id).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("GetUserByID: user not found for id:", id)
			return nil, err
		}
		log.Println("GetUserByID: database error:", err)
		return nil, err
	}
	log.Println("GetUserByID: user found:", user)
	return &user, nil
}

func (s *Storage) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := s.db.QueryRow("SELECT id, email, firstName, lastName, password, createdAt FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName,&user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("GetUserByEmail: user not found for email:", email)
			return nil, err
		}
		log.Println("GetUserByEmail: database error:", err)
		return nil, err
	}
	log.Println("GetUserByEmail: user found:", user)
	return &user, nil
}


func (s *Storage) GetTask(id string) (*model.Task, error) {
	var task model.Task
	err := s.db.QueryRow("SELECT id, name, status, projectId, assignedToID, createdAt FROM tasks WHERE id = ?", id).Scan(&task.ID, &task.Name, &task.Status, &task.ProjectID, &task.AssignedToID, &task.CreatedAt)
	return &task, err
}

func (s *Storage) GetProject(id string) (*model.Project, error) {
	var project model.Project
	err := s.db.QueryRow("SELECT id, name, createdAt FROM projects WHERE id = ?", id).Scan(&project.ID, &project.Name, &project.CreatedAt)
	return &project, err
}

func (s *Storage) DeleteProject(id string) error {
	result, err := s.db.Exec("DELETE FROM projects WHERE id = ?", id)
	if err != nil {
			log.Printf("DeleteProject: error executing query: %v", err)
			return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
			log.Printf("DeleteProject: error getting affected rows: %v", err)
			return err
	}

	if rowsAffected == 0 {
			log.Printf("DeleteProject: no rows affected, project id %s not found", id)
			return sql.ErrNoRows
	}

	log.Printf("DeleteProject: project id %s deleted successfully", id)
	return nil
}

func (s *Storage) DeleteTask(id string) error {
	_, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
			log.Printf("DeleteTask: error executing query: %v", err)
			return err
	}
	log.Printf("DeleteTask: task with id %s deleted successfully", id)
	return nil
}

func (s *Storage) UpdateTask(t *model.Task) (*model.Task, error) {
	_, err := s.db.Exec("UPDATE tasks SET name = ?, status = ?, projectId = ?, assignedToID = ? WHERE id = ?", t.Name, t.Status, t.ProjectID, t.AssignedToID, t.ID)
	if err != nil {
		log.Println("UpdateTask: error executing query:", err)
		return nil, err
	}
	return t, nil
}

func (s *Storage) UpdateProject(p *model.Project) (*model.Project, error) {
	_, err := s.db.Exec("UPDATE projects SET name = ? WHERE id = ?", p.Name, p.ID)
	if err != nil {
		log.Println("UpdateProject: error executing query:", err)
		return nil, err
	}
	return p, nil
}