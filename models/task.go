package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	AssignedTo  []string           `bson:"assignedTo,omitempty"` // ✅ Now stores user IDs as strings
	AssignedBy  string             `bson:"assignedBy,omitempty"` // ✅ Now stores assigner ID as string
	Status      string             `bson:"status"`
	Completed   bool               `bson:"completed"`
}
