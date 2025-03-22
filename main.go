package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data.db"
	}
	adminKey := os.Getenv("ADMIN_KEY")
	if adminKey == "" {
		log.Fatal("ADMIN_KEY environment variable is required")
	}

	// get port from env or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db := InitDB(dbPath)
	defer db.Close()

	store := &Storage{DB: db}
	handler := &Handler{Store: store, AdminKey: adminKey}

	http.HandleFunc("/shorten", withDB(db, BasicAuth(handler.CreateOrUpdateShortLink)))
	http.HandleFunc("/shorten/", withDB(db, BasicAuth(handler.GetShortLink)))
	http.HandleFunc("/delete/", withDB(db, BasicAuth(handler.DeleteShortLink)))
	http.HandleFunc("/admin/create-user", handler.CreateUser)
	http.HandleFunc("/r/", handler.Redirect)

	log.Println("Server started at :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// withDB injects the database into the request context
func withDB(db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "db", db)
		next(w, r.WithContext(ctx))
	}
}
