package main

import (
	"log"
	"os"

	"task-management/config"
	"task-management/controllers"
	"task-management/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println(" Warning: No .env file found. Using system environment variables.")
	}

	// Connect to MongoDB
	log.Println("🔗 Connecting to MongoDB...")
	if err := config.ConnectDB(); err != nil {
		log.Fatalf("MongoDB Connection Failed: %v", err)
	}
	log.Println("Successfully connected to MongoDB.")

	// Initialize collections
	log.Println("🔹 Initializing collections...")
	controllers.InitUserCollection()
	controllers.InitTaskCollection()

	// Create Gin Router
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000", "https://your-frontend-domain.vercel.app"}, // Add your deployed frontend URL here
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Define Routes

	// Authentication Routes
	r.POST("/auth/register", controllers.RegisterUser)
	r.POST("/auth/login", controllers.LoginUser)

	// Task Routes
	r.GET("/tasks", controllers.GetTasks)
	r.POST("/tasks", controllers.CreateTask)
	r.PUT("/tasks/:id/assign", controllers.AssignTask)
	r.PUT("/tasks/:id/status", controllers.UpdateTaskStatus)
	r.GET("/tasks/completed", controllers.GetCompletedTasks)
	r.GET("/tasks/in-progress", controllers.GetInProgressTasks)
	r.GET("/tasks/pending", controllers.GetPendingTasks)

	// WebSocket Route
	r.GET("/ws", websocket.HandleWebSocket)

	// Get PORT from Environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 if not set
	}

	// Start the Server
	log.Printf(" Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf(" Failed to start server: %v", err)
	}
}
