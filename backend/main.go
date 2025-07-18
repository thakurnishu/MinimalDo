package main

import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	db *sql.DB
)


func main() {
	cfg := loadConfig()

	// Otel init
	cleanup, err := InitTelemetry(cfg)
	if err != nil {
		slog.Error("Failed to initialize telemetry", "error", err)
	}

	defer cleanup()

	db := setupDB(cfg)
	defer db.Close()
	server := &Server{
		db: db,
	}

	router := gin.Default()

	// CORS setup
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{cfg.FrontendURL},
		AllowHeaders: []string{"X-Requested-With", "Content-Type", "Authorization"},
		AllowMethods: []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"},
		ExposeHeaders: []string{"Content-Length"},
	}))
	router.Use(otelgin.Middleware(cfg.ServiceName))

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
	
	slog.Info("server is listening", "port", cfg.Port)
	router.Run(":"+cfg.Port)
}
