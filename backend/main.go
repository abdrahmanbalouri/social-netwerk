package main

import (
	"log"
	"net/http"

	"social-network/internal/api"
	"social-network/internal/repository"
	// "social-network/pkg/middlewares"
	// "social-network/pkg/ratelimiter"
)

func main() {
	db, err := repository.OpenDb()
	if err != nil {
		log.Fatal("Error: ", err)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			db.Close()
			log.Fatal("Error: ", err)
		}
	}()

	if err := repository.ApplyMigrations(db); err != nil {
		panic("Migration failed: " + err.Error())
	}

	baseHandler := api.Routes(db)

	server := &http.Server{
		Addr:    ":8080",
		Handler: baseHandler,
	}
	log.Println("http://localhost:8080/")
	err = server.ListenAndServe()
	if err != nil {
		log.Println("Error in starting of server:", err)
		db.Close()
		return
	}
}
