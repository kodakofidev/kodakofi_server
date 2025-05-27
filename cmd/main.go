package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kodakofidev/kodakofi_server/internal/handlers"
	"github.com/kodakofidev/kodakofi_server/internal/models"
	"github.com/kodakofidev/kodakofi_server/internal/routes"
	"github.com/kodakofidev/kodakofi_server/pkg"
)

var dbpool *pgxpool.Pool

func main() {
	// Load environment variables early
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize OAuth providers before setting up routes
	handlers.InitAuth()

	db, err := pkg.Posql()
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer db.Close()

	log.Println("DB connected successfully")

	router := routes.InitRouter(db)

	router.Static("/public/profile-images", "./public/profile-images")

	router.Static("/public/product-image", "./public/product-image")

	router.GET("/ping", func(c *gin.Context) {
		responder := models.NewResponse(c)
		ctx := c.Request.Context()
		if err := db.Ping(ctx); err != nil {
			responder.InternalServerError("Database not responding", err.Error())
			return
		}
		responder.Success("pong", nil)
	})

	srv := pkg.Server(router)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
