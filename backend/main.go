package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	// Database connection
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "todoapp")
	
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Initialize database
	if err := initDB(db); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	server := &Server{db: db}

	// CORS setup
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"X-Requested-With", "Content-Type", "Authorization"},
		AllowMethods: []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	// Setup routes
	api := router.Group("/api")
	api.GET("/todos", server.getTodos)
	api.POST("/todos", server.createTodo)
	api.PUT("/todos/:id", server.updateTodo)
	api.DELETE("/todos/:id", server.deleteTodo)
	api.GET("/health", server.healthCheck)

	
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	router.Run(":"+port)
}

func initDB(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		completed BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE OR REPLACE FUNCTION update_updated_at_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = CURRENT_TIMESTAMP;
		RETURN NEW;
	END;
	$$ language 'plpgsql';

	DROP TRIGGER IF EXISTS update_todos_updated_at ON todos;
	CREATE TRIGGER update_todos_updated_at
		BEFORE UPDATE ON todos
		FOR EACH ROW
		EXECUTE FUNCTION update_updated_at_column();
	`
	
	_, err := db.Exec(query)
	return err
}

func (s *Server) getTodos(c *gin.Context){
	rows, err := s.db.Query(`
		SELECT id, title, description, completed, created_at, updated_at 
		FROM todos 
		ORDER BY created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		todos = append(todos, t)
	}
	
	c.JSON(http.StatusOK, todos)
}


func (s *Server) createTodo(c *gin.Context){
	var t Todo
	if err := c.BindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		INSERT INTO todos (title, description, completed) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at
	`
	
	err := s.db.QueryRow(query, t.Title, t.Description, t.Completed).
		Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK,t)	
}

func (s *Server) updateTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invaild ID"})
		return
	}

	var t Todo
	if err := c.BindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	query := `
		UPDATE todos 
		SET title = $1, description = $2, completed = $3
		WHERE id = $4
		RETURNING id, title, description, completed, created_at, updated_at
	`
	
	err = s.db.QueryRow(query, t.Title, t.Description, t.Completed, id).
		Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, t)
}


func (s *Server) deleteTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invaild ID"})
	}

	result, err := s.db.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not fount"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (s *Server) healthCheck(c *gin.Context) {
	response := map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
