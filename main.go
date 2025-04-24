package main

import (
	"context"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tenteedee/go-mongo/controllers"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var client *mongo.Client

func main() {
	client = getSession()
	defer disconnectMongoClient()

	router := httprouter.New()
	todoController := controllers.NewTodoController(client)

	router.GET("/todo/:id", todoController.GetTodoById)
	router.POST("/todo", todoController.CreateTodo)
	router.PUT("/todo/:id", todoController.UpdateTodo)
	router.DELETE("/todo/:id", todoController.DeleteTodo)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Server error:", err)
	}
}

func getSession() *mongo.Client {
	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	log.Println("Connected to MongoDB")

	return client
}

func disconnectMongoClient() {
	if client != nil {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Fatal("Error disconnecting MongoDB:", err)
		}
		log.Println("MongoDB connection closed")
	}
}
