package main

import (
	"context"
	"log"
	"os"
	"time"

	"detection-api/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load env (gunakan godotenv jika perlu)
if os.Getenv("APP_ENV") != "production" {
    err := godotenv.Load()
    if err != nil {
        log.Println("⚠️ Warning: .env file not found, using system environment variables")
    }
}


	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("❌ MongoDB connection error:", err)
	}
	defer client.Disconnect(ctx)

	// Ping database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("❌ MongoDB ping failed:", err)
	}
	log.Println("✅ MongoDB connected")

	// Setup Gin
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:    []string{("*")},
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:   []string{"Content-Length"},
	}))



	// Initialize handler
	collection := client.Database("go-backend").Collection("detections")
	detectionHandler := &handlers.DetectionHandler{Collection: collection}

	// Routes
	api := r.Group("/api")
	{
		api.POST("/detections", detectionHandler.CreateDetection)
	}

	// Start server
	log.Printf("🚀 Server running on port %s", port)
	r.Run(":" + port)
}