package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	db *sql.DB
)

func init() {
	loadConfig()
}

func main() {
	// Otel 
	cleanup := initTracer()
	defer func() {
		if err := cleanup(context.Background()); err != nil {
			log.Fatalf("Failed to shutdown tracer: %v", err)
		}
	}()

	db := setupDB()
	defer db.Close()
	server := &Server{
		db: db,
	}

	router := gin.Default()

	// CORS setup
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{frontendURL},
		AllowHeaders: []string{"X-Requested-With", "Content-Type", "Authorization"},
		AllowMethods: []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"},
		ExposeHeaders: []string{"Content-Length"},
	}))
	router.Use(otelgin.Middleware(serviceName))

	// Setup routes
	api := router.Group("/api")
	{
		api.GET("/todos", server.getTodos)
		api.POST("/todos", server.createTodo)
		api.PUT("/todos/:id", server.updateTodo)
		api.DELETE("/todos/:id", server.deleteTodo)
		api.GET("/health", server.healthCheck)
		api.GET("/todos/by-date", server.getTodosByDate)
	}
	
	log.Printf("Server starting on port %s", port)
	router.Run(":"+port)
}
