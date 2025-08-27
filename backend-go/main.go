package main

import (
	"log"
	"os"

	"github.com/aliakkas006/backend-go/db"
	"github.com/aliakkas006/backend-go/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/aliakkas006/backend-go/controllers"
	"github.com/aliakkas006/backend-go/kafka"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	brokers := os.Getenv("BOOTSTRAP_SERVERS")
	topic := os.Getenv("TOPIC_NAME")

	p := kafka.NewProducer(brokers, topic)
	defer p.Close()
	controllers.InitProducer(p)

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}

	err = db.Connect()
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v", err)
	}
	defer db.Close()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	// Add a simple health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	routes.RegisterTodoRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Starting Go backend on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
