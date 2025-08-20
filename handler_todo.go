package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"todo-app/internal/database"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiConfig *ApiConfig) getAllTodos(w http.ResponseWriter, r *http.Request) {

	todos, err := apiConfig.db.GetAllTodos(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error retrieving todos: %v", err))
		return
	}
	if len(todos) == 0 {
		respondWithError(w, http.StatusNotFound, "No todos found")
		return
	}
	var todoList []Todo

	for _, todo := range todos {
		todoList = append(todoList, databaseObjectToJsonObject(todo))
	}
	responseWithJson(w, http.StatusOK, todoList)
}

func (apiConfig *ApiConfig) addTodoHandler(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	decoder := json.NewDecoder(r.Body)
	payload := Payload{}
	err := decoder.Decode(&payload)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to decode payload")
		return
	}

	if payload.Title == "" || payload.Description == "" {
		respondWithError(w, http.StatusBadRequest, "Missing params")
		return
	}

	todo, err := apiConfig.db.CreateTodo(r.Context(), database.CreateTodoParams{
		ID:          uuid.New(),
		Title:       payload.Title,
		Description: payload.Description,
		Status:      false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to save the to-do task")
		return
	}
	responseWithJson(w, http.StatusOK, databaseObjectToJsonObject(todo))
}

func (apiConfig ApiConfig) handleUpdateTodoStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "id missing from the url")
		return
	}
	todoId, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid id")
		return
	}
	res, err := apiConfig.db.MarkTodoAsDone(r.Context(), database.MarkTodoAsDoneParams{
		Status:    true,
		UpdatedAt: time.Now(),
		ID:        todoId,
	})

	if err != nil {
		respondWithError(w, http.StatusNotFound, "No todo task found with this id")
		return
	}
	responseWithJson(w, http.StatusOK, databaseObjectToJsonObject(res))
}

func (apiConfig ApiConfig) deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	todoId, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid to-do ID")
		return
	}

	databaseTodo, err := apiConfig.db.GetTodoByID(r.Context(), todoId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No to-do task found with this id")
		return
	}

	err = apiConfig.db.DeleteAtodo(r.Context(), databaseTodo.ID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid to-do ID")
		return
	}
	type Success struct {
		Message string `json:"message"`
	}
	responseWithJson(w, http.StatusOK, Success{
		Message: "To-do task deleted successfully",
	})
}

func (apiConfig ApiConfig) filterTodosHander(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status == "" {
		respondWithError(w, http.StatusBadRequest, "Status query parameter is required")
		return
	}

	var statusBool bool
	switch status {
	case "true":
		statusBool = true
	case "false":
		statusBool = false
	default:
		respondWithError(w, http.StatusBadRequest, "Invalid status value")
		return
	}

	todos, err := apiConfig.db.FilterTodos(r.Context(), statusBool)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error retrieving todos: %v", err))
		return
	}

	if len(todos) == 0 {
		respondWithError(w, http.StatusNotFound, "No todos found with the specified status")
		return
	}

	var todoList []Todo
	for _, todo := range todos {
		todoList = append(todoList, databaseObjectToJsonObject(todo))
	}
	responseWithJson(w, http.StatusOK, todoList)

}
