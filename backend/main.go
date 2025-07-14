package main

import (
	"context"
	"log"

	_ "github.com/lib/pq"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")
	log.Printf("CORS: Allowing origin: %s", frontendURL)

	// Otel 
	cleanup := initTracer()
	defer cleanup(context.Background())

	db := setupDB()
	defer db.Close()
	server := &Server{db: db}


	// CORS setup
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{frontendURL},
		AllowHeaders: []string{"X-Requested-With", "Content-Type", "Authorization"},
		AllowMethods: []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"},
		ExposeHeaders: []string{"Content-Length"},
	}))
	router.Use(otelgin.Middleware(serviceName))

	// Setup routes
	api := router.Group("/api")

//	api.Use(otelgin.Middleware(serviceName))
//
//	api.Use(cors.New(cors.Config{
//		AllowOrigins: []string{frontendURL},
//		AllowHeaders: []string{"X-Requested-With", "Content-Type", "Authorization"},
//		AllowMethods: []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"},
//		ExposeHeaders: []string{"Content-Length"},
//	}))
	{
		api.GET("/todos", server.getTodos)
		api.POST("/todos", server.createTodo)
		api.PUT("/todos/:id", server.updateTodo)
		api.DELETE("/todos/:id", server.deleteTodo)
		api.GET("/health", server.healthCheck)
	}
	
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	router.Run(":"+port)
}
