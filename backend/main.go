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

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		db.Close()
	// 		log.Fatal("Error: ", err)
	// 	}
	// }()

	if err := repository.ApplyMigrations(db); err != nil {
		panic("Migration failed: " + err.Error())
	}

	baseHandler := api.Routes()

	// Wrap the API routes with CORS
	handler := enableCORS(baseHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	/* var rr []string
	row, _ := repository.Db.Query(`select id from users where id !=  '847334d1-e080-4536-bdb4-256708383ef0'`)

	for row.Next() {
		var id string
		row.Scan(&id)
		rr = append(rr, id)
	}

	for i := 0; i < len(rr); i++ {
		_, err := repository.Db.Exec("insert into  follow_requests (user_id , follower_id) values (?,?)", "847334d1-e080-4536-bdb4-256708383ef0", rr[i])
		if err != nil {
			fmt.Println("erfref", err)
			continue
		}
	} */

	log.Println("http://localhost:8080/")
	err = server.ListenAndServe()
	if err != nil {
		log.Println("Error in starting of server:", err)
		db.Close()
		return
	}
}
