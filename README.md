# Go Rest Project Manager

## Brief Description

Go Rest Project Manager is a REST API-based project manager application developed using the Go programming language. This application allows users to easily manage their projects and tasks, such as creating, reading, updating, and deleting (CRUD) project and task data.

## Technology Stack

- **Programming Language**: Go (Golang)
- **Database**: MySQL
- **API Framework**: net/http (Go standard library)
- **Containerization**: Docker
- **Dependency Management**: Go Modules

## Features

- **User Management**: User registration and login.
- **Project Management**: CRUD projects.
- **Task Management**: CRUD tasks.
- **Authentication**: JWT-based authentication for API security.
- **Input Validation**: Validation of user input data.
- **API Documentation**: Endpoint documentation using Postman Collection.
- **Testing**: Includes comprehensive testing to ensure reliability and correctness of the API endpoints.

## How to Run the Program

### Prerequisites

Ensure you have installed:

- [Docker](https://www.docker.com/products/docker-desktop)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Go](https://golang.org/dl/)

### Steps

1. **Clone Repository**

```bash
git clone https://github.com/DaffaJatmiko/go-rest-project-manager.git
cd go-rest-project-manager
```

2. **Configure Environment Variables**

Create a `.env` file in the root directory and add the following configuration:

```bash
DB_USER=root
DB_PASSWORD=mysqldatabase123
DB_HOST=mysql
DB_PORT=3306
DB_NAME=goprojectmanager
JWT_SECRET=randomjwtsecret
```

3. **Build and Run Docker Containers**

```bash
docker-compose up --build
```

4. **Accessing the API**

The application will run at `http://localhost:3000`. You can use Postman or other tools to access the API endpoints.

## Testing

This project includes testing to ensure the reliability and correctness of the API endpoints. You can find the testing code in the service folder. To run the tests, follow these steps:

1. Navigate to the service folder in your terminal.
2. Run the following command to execute the tests:

```bash
go test
```
