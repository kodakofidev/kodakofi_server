package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	"github.com/kodakofidev/kodakofi_server/internal/models"
	"github.com/kodakofidev/kodakofi_server/internal/routes"
	"github.com/kodakofidev/kodakofi_server/pkg"
)

var dbpool *pgxpool.Pool

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	db, err := pkg.Posql()
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer db.Close()

	log.Println("DB connected successfully")

	router := routes.InitRouter(db)

	router.Static("/img", "./public/img")

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
