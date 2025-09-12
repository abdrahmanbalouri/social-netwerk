package main

import (
	"database/sql"
	"log"
	"net/http"

	"social-network/internal/api"
	"social-network/internal/repository"
	// "social-network/pkg/middlewares"
	// "social-network/pkg/ratelimiter"
)

var Db *sql.DB

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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
	hand := enableCORS(baseHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: hand,
	}
	log.Println("http://localhost:8080/")
	err = server.ListenAndServe()
	if err != nil {
		log.Println("Error in starting of server:", err)
		db.Close()
		return
	}
}
