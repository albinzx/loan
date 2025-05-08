package main

import (
	"log"

	"github.com/albinzx/loan/pkg/sql"
	"github.com/albinzx/loan/pkg/sql/mysql"
	"github.com/albinzx/loan/repository"
	"github.com/albinzx/loan/service"
	"github.com/albinzx/loan/transport/http"
)

func main() {
	db, err := sql.DB(&mysql.DataSource{
		Host:      "127.0.0.1",
		Port:      "3306",
		User:      "dev",
		Password:  "dev",
		Database:  "amartha_loan",
		ParseTime: true,
	})

	if err != nil {
		log.Printf("error while connecting to DB, %v", err)
		return
	}

	repo := repository.New(db)
	svc := service.New(repo)
	trp := http.New(svc)

	trp.Serve("127.0.0.1:8080", "/v1")
}
