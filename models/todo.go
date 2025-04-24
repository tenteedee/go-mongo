package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Todo struct {
	ID        bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string        `json:"title" bson:"title"`
	Completed bool          `json:"completed" bson:"completed" default:"false"`
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}
