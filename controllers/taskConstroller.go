package controllers

import (
	"context"
	"log"
	"net/http"
	"task-management/config"
	"task-management/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var taskCollection *mongo.Collection

func InitTaskCollection() {
	if config.DB == nil {
		log.Fatal("❌ Database is not initialized. Call ConnectDB() first.")
	}
	taskCollection = config.DB.Collection("tasks")
	log.Println("✅ Task collection initialized")
}

// Create Task
func CreateTask(c *gin.Context) {
	if taskCollection == nil {
		log.Println("❌ Task collection is nil! Ensure InitTaskCollection() is called.")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
		return
	}

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task data"})
		return
	}

	//  Assign a unique MongoDB ObjectID
	task.ID = primitive.NewObjectID()
	task.Status = "Pending"
	task.Completed = false

	_, err := taskCollection.InsertOne(context.TODO(), task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Task created", "task": task})
}

// Get All Tasks
func GetTasks(c *gin.Context) {
	if taskCollection == nil {
		log.Println("❌ Task collection is nil! Ensure InitTaskCollection() is called.")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
		return
	}

	//  Fetch all tasks from MongoDB
	cursor, err := taskCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println("❌ MongoDB Query Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tasks"})
		return
	}

	var tasks []models.Task
	if err := cursor.All(context.TODO(), &tasks); err != nil {
		log.Println("❌ Decoding Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding tasks"})
		return
	}

	// Log raw tasks for debugging
	log.Println("✅ Retrieved Tasks:", tasks)

	// If no tasks are found, return an empty array
	if len(tasks) == 0 {
		c.JSON(http.StatusOK, gin.H{"tasks": []models.Task{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// Assign Task to Multiple Users
func AssignTask(c *gin.Context) {
	taskID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var request struct {
		AssignedTo []string `json:"assignedTo"` // ✅ Accepts multiple assignees as strings
		AssignedBy string   `json:"assignedBy"` // ✅ Accepts the assigner ID as a string
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	update := bson.M{"$set": bson.M{
		"assignedTo": request.AssignedTo,
		"assignedBy": request.AssignedBy,
		"status":     "Pending",
	}}

	_, err = taskCollection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Task assigned",
		"assignedTo": request.AssignedTo,
		"assignedBy": request.AssignedBy,
	})
}

// Update Task Status
func UpdateTaskStatus(c *gin.Context) {
	taskID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var request struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	//  Validate allowed statuses
	if request.Status != "Pending" && request.Status != "In Progress" && request.Status != "Completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	//  Automatically mark completed if status is "Completed"
	completedFlag := request.Status == "Completed"

	update := bson.M{"$set": bson.M{
		"status":    request.Status,
		"completed": completedFlag,
	}}

	_, err = taskCollection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Task status updated",
		"status":    request.Status,
		"completed": completedFlag,
	})
}

// Get Tasks by Status

// GetCompletedTasks fetches all tasks with status "Completed"
func GetCompletedTasks(c *gin.Context) {
	cursor, err := taskCollection.Find(context.TODO(), bson.M{"status": bson.M{"$regex": "^Completed$", "$options": "i"}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch completed tasks"})
		return
	}

	var tasks []models.Task
	if err := cursor.All(context.TODO(), &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing completed tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func GetInProgressTasks(c *gin.Context) {
	cursor, err := taskCollection.Find(context.TODO(), bson.M{"status": bson.M{"$regex": "^In Progress$", "$options": "i"}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch in-progress tasks"})
		return
	}

	var tasks []models.Task
	if err := cursor.All(context.TODO(), &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing in-progress tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func GetPendingTasks(c *gin.Context) {
	cursor, err := taskCollection.Find(context.TODO(), bson.M{"status": bson.M{"$regex": "^Pending$", "$options": "i"}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending tasks"})
		return
	}

	var tasks []models.Task
	if err := cursor.All(context.TODO(), &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing pending tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}
