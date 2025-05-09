package main

import (
	"log"

	"github.com/albinzx/loan/pkg/config/viper"
	"github.com/albinzx/loan/pkg/mailer"
	"github.com/albinzx/loan/pkg/sql"
	"github.com/albinzx/loan/pkg/sql/mysql"
	"github.com/albinzx/loan/repository"
	"github.com/albinzx/loan/service"
	"github.com/albinzx/loan/transport/http"
)

func main() {
	cfg, err := viper.New(".", "config", "json")

	if err != nil {
		log.Printf("error while reading config, %v", err)
		return
	}

	db, err := sql.DB(&mysql.DataSource{
		Host:      cfg.GetString("db.host"),
		Port:      cfg.GetString("db.port"),
		User:      cfg.GetString("db.user"),
		Password:  cfg.GetString("db.password"),
		Database:  cfg.GetString("db.name"),
		ParseTime: true,
	})

	if err != nil {
		log.Printf("error while connecting to DB, %v", err)
		return
	}

	mail := mailer.New(cfg)
	repo := repository.New(db)
	svc := service.New(repo, mail)
	trp := http.New(svc)

	trp.Serve(
		cfg.GetString("http.address"),
		cfg.GetString("http.path"),
	)
}
