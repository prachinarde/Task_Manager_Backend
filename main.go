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

	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Warning: No .env file found. Using system environment variables.")
	}

	log.Println("🔗 Connecting to MongoDB...")
	if err := config.ConnectDB(); err != nil {
		log.Fatalf("❌ MongoDB Connection Failed: %v", err)
	}
	log.Println("✅ Successfully connected to MongoDB.")

	log.Println("🔹 Initializing collections...")
	controllers.InitUserCollection()
	controllers.InitTaskCollection()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Authentication Routes
	r.POST("/auth/register", controllers.RegisterUser)
	r.POST("/auth/login", controllers.LoginUser)

	// Task Routes
	r.GET("/tasks", controllers.GetTasks)
	r.POST("/tasks", controllers.CreateTask)
	r.PUT("/tasks/:id/assign", controllers.AssignTask)          // Assign Task
	r.PUT("/tasks/:id/status", controllers.UpdateTaskStatus)    // Update Task Status
	r.GET("/tasks/completed", controllers.GetCompletedTasks)    // Get all completed tasks
	r.GET("/tasks/in-progress", controllers.GetInProgressTasks) // Get all in-progress tasks
	r.GET("/tasks/pending", controllers.GetPendingTasks)        // Get all pending tasks

	// WebSocket Route
	r.GET("/ws", websocket.HandleWebSocket)

	// Define Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080
	}

	log.Println("Server running on port", port)
	r.Run(":" + port)
}
