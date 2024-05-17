package main

import (
	"log"

	"github.com/DaffaJatmiko/go-rest-project-manager/cmd/api"
	"github.com/DaffaJatmiko/go-rest-project-manager/config"
	"github.com/DaffaJatmiko/go-rest-project-manager/db"
	"github.com/DaffaJatmiko/go-rest-project-manager/repository"
	"github.com/go-sql-driver/mysql"
)

func main() {
	cfg := mysql.Config{
		User: 								config.Envs.DBUser,
		Passwd: 							config.Envs.DBPassword,
		Addr: 								config.Envs.DBAddress,
		DBName: 							config.Envs.DBName,
		Net: 									"tcp",
		AllowNativePasswords: true,
		ParseTime: 						true,
	}

	sqlStorage := db.NewMysqlStorage(cfg)

	db, err := sqlStorage.Init()
	if err != nil {
		log.Fatal(err)
	}

	store := repository.NewStore(db)
	api := api.NewAPIServer(":3000", store)
	api.Serve()
}