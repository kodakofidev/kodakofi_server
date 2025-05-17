package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	"github.com/kodakofidev/kodakofi_server/internal/routes"
	"github.com/kodakofidev/kodakofi_server/pkg"
)

var dbpool *pgxpool.Pool

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	var err error
	dbpool, err = pkg.Posql()
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer dbpool.Close()

	log.Println("DB connected successfully")

	http.HandleFunc("/ping", pingHandler)
	router := routes.InitRouter(dbpool)

	// log.Println("Server listening on :8080")
	// if err := (":8000", nil); err != nil {
	// 	log.Fatal("HTTP server failed:", err)
	// }
	router.Run("localhost:8000")
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := dbpool.Ping(ctx); err != nil {
		http.Error(w, "Database not responding", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
