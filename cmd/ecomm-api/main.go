package main

import (
	"log"

	"github.com/m21power/ecomm/db"
	"github.com/m21power/ecomm/ecomm-api/handler"
	"github.com/m21power/ecomm/ecomm-api/server"
	"github.com/m21power/ecomm/ecomm-api/storer"
)

func main() {
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	defer db.Close()
	log.Printf("Connected to database")
	st := storer.NewMySQLStorer(db.GetDB())
	server := server.NewServer(st)
	h := handler.NewHandler(server)
	handler.RegisterRoutes(h)
	log.Printf("Starting server on :8080")
	err = handler.Start(":8080")
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
