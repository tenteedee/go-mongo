package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tenteedee/go-mongo/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TodoController struct {
	session *mongo.Client
}

func NewTodoController(session *mongo.Client) *TodoController {
	return &TodoController{
		session: session,
	}
}

func (tc *TodoController) GetTodoById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		NotFound(w, r, p, fmt.Errorf("id is not valid"))
		return
	}

	todo := models.Todo{}
	collection := tc.session.Database("todo").Collection("todos")
	err = collection.FindOne(r.Context(), bson.M{"_id": objectId}).Decode(&todo)
	if err != nil {
		NotFound(w, r, p, err)
		return
	}

	todoJson, err := json.Marshal(todo)
	if err != nil {
		InternalServerError(w, r, p, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(todoJson)
}

func (tc *TodoController) CreateTodo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	todo := models.CreateTodoRequest{}

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		InternalServerError(w, r, p, err)
		return
	}

	collection := tc.session.Database("todo").Collection("todos")
	result, err := collection.InsertOne(r.Context(), todo)
	if err != nil {
		InternalServerError(w, r, p, err)
		return
	}

	todoRes := models.Todo{
		ID:        (result.InsertedID).(bson.ObjectID),
		Title:     todo.Title,
		Completed: false,
	}

	todoJson, err := json.Marshal(todoRes)
	if err != nil {
		InternalServerError(w, r, p, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(todoJson)
}

func (tc *TodoController) DeleteTodo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		NotFound(w, r, p, fmt.Errorf("id is not valid"))
		return
	}

	collection := tc.session.Database("todo").Collection("todos")
	_, err = collection.DeleteOne(r.Context(), bson.M{"_id": objectId})
	if err != nil {
		InternalServerError(w, r, p, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (tc *TodoController) UpdateTodo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		NotFound(w, r, p, fmt.Errorf("id is not valid"))
		return
	}

	var updatedTodo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		InternalServerError(w, r, p, err)
		return
	}

	updateDoc := bson.M{}
	if updatedTodo.Title != "" {
		updateDoc["title"] = updatedTodo.Title
	}

	if updatedTodo.Completed {
		updateDoc["completed"] = updatedTodo.Completed
	}

	if len(updateDoc) == 0 {
		BadRequest(w, r, p, fmt.Errorf("no valid fields to update"))
		return
	}

	updatedTodo.ID = objectId

	collection := tc.session.Database("todo").Collection("todos")
	updateResult, err := collection.UpdateOne(
		r.Context(),
		bson.M{"_id": objectId},
		bson.M{
			"$set": updateDoc,
		},
	)
	if err != nil {
		InternalServerError(w, r, p, err)
		return
	}

	if updateResult.MatchedCount == 0 {
		NotFound(w, r, p, fmt.Errorf("todo with id %s not found", id))
		return
	}

	updatedTodoJson, err := json.Marshal(updatedTodo)
	if err != nil {
		InternalServerError(w, r, p, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(updatedTodoJson)
}
