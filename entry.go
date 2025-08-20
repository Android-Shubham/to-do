package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	_ "github.com/lib/pq"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"todo-app/internal/database"
)

type ApiConfig struct {
	db *database.Queries
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("No PORT environment variable detected, using default port 8080")
	}else{
		fmt.Printf("Using PORT: %s\n", port)
	}


	db_url := os.Getenv("DB_URL")
	if(db_url == ""){
		fmt.Println("No DB_URL environment variable detected, using default database URL")
	}

	conn,err := sql.Open("postgres", db_url)
	if err != nil {
		fmt.Printf("Error connecting to the database: %v\n", err)
	}

	queries := database.New(conn)

	apiConfig := ApiConfig{
		db: queries,
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: [] string{"*"},
		AllowedMethods: [] string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: [] string{"*"},
		AllowCredentials: false,
		MaxAge: 300, // Maximum age in seconds for preflight requests
	}))

	v1Router := chi.NewRouter()
	//endpoints here
	v1Router.Get("/healthz", handleSuccess)
	v1Router.Get("/error", handleError)
	v1Router.Get("/todos", apiConfig.getAllTodos)
	v1Router.Post("/todos",apiConfig.addTodoHandler)
	v1Router.Put("/todos/{id}",apiConfig.handleUpdateTodoStatus)
	v1Router.Delete("/todos/{id}",apiConfig.deleteTodoHandler)
	v1Router.Get("/todos/status",apiConfig.filterTodosHander)

	router.Mount("/v1", v1Router)
	srv := &http.Server{
		Handler: router,
		Addr:":"+port,
	}
	err= srv.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
