package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	api "social-network/app/router"
	middlewares "social-network/middleware"
	"social-network/pkg/db/sqlite"
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
	fmt.Println("4545")
	db, err := sqlite.OpenDb()
	if err != nil {
		log.Fatal("Error: ", err)
		return
	}
	if err := sqlite.ApplyMigrations(db); err != nil {
		panic("Migration failed: " + err.Error())
	}

	baseHandler := api.Routes()

	// Wrap the API routes with CORS
	handler := enableCORS(middlewares.SessionMiddleware(sqlite.Db, baseHandler))

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Println("http://localhost:8080/")
	err = server.ListenAndServe()
	if err != nil {
		log.Println("Error in starting of server:", err)
		db.Close()
		return
	}
}
