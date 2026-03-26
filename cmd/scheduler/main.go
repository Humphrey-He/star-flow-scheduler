package main

import (
    "log"
    "net/http"

    "starflow-scheduler/internal/config"
    "starflow-scheduler/internal/db"
    "starflow-scheduler/internal/httpapi"
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
