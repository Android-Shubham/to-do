package main

import (
	"time"
	"todo-app/internal/database"

	"github.com/google/uuid"
)

type Todo struct {
	ID uuid.UUID `json:"id"`
	Title string `json:"title"`
	Description string `json:"description,omitempty"`
	Status bool `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func databaseObjectToJsonObject(todo database.Todo) Todo {
	return Todo{
		ID: todo.ID,
		Title: todo.Title,
		Description: todo.Description,
		Status: todo.Status,
		CreatedAt: todo.CreatedAt,
		UpdatedAt: todo.UpdatedAt,
	}
}