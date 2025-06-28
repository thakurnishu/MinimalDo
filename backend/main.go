package main

import (
    "log"
    "net/http"
    "os"

		"github.com/thakurnishu/MinimalDo/database"
		"github.com/thakurnishu/MinimalDo/handlers"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
)

func main() {
    // Initialize database
    database.InitDB()
    defer database.DB.Close()

    // Setup routes
    r := mux.NewRouter()
    
    // API routes
    api := r.PathPrefix("/api").Subrouter()
    api.HandleFunc("/todos", handlers.GetTodos).Methods("GET")
    api.HandleFunc("/todos", handlers.CreateTodo).Methods("POST")
    api.HandleFunc("/todos/{id}", handlers.UpdateTodo).Methods("PUT")
    api.HandleFunc("/todos/{id}", handlers.DeleteTodo).Methods("DELETE")

    // Health check
    r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    }).Methods("GET")

    // Setup CORS
    c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"}, // In production, specify your frontend URL
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"*"},
    })

    handler := c.Handler(r)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, handler))
}
