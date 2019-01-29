package main

import (
	"flag"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/mcornut/go-rest-api/controllers"
	"github.com/mcornut/go-rest-api/routes"
	"github.com/mcornut/go-rest-api/utils/database"
)

//go:generate mockery -name DocumentController -dir ./controllers
//go:generate mockery -name DocumentRepository -dir ./repositories

func formatLog() {
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
}

func main() {

	formatLog()

	configPath := flag.String("config", "local.toml", "config path")
	flag.Parse()

	config, err := ConfigFromFile(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(config.DB.Username, config.DB.Password, config.DB.Name, config.DB.Host, config.DB.Port)
	if err != nil {
		log.Fatal(err)
	}

	documentController := controllers.NewDocumentController(db)

	mux := http.NewServeMux()
	routes.CreateRoutes(mux, documentController)

	addr := "localhost:8080"
	log.Infof("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
