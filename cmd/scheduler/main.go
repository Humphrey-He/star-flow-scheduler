package main

import (
	"log"
	"net/http"

	"github.com/Humphrey-He/star-flow-scheduler/internal/config"
	"github.com/Humphrey-He/star-flow-scheduler/internal/db"
	"github.com/Humphrey-He/star-flow-scheduler/internal/httpapi"
)

func main() {
	cfg := config.Load()

	database, err := db.Open(cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("db open failed: %v", err)
	}
	defer database.Close()

	server := httpapi.NewServer(database)

	log.Printf("http server listening on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, server.Routes()); err != nil {
		log.Fatalf("http server failed: %v", err)
	}
}
