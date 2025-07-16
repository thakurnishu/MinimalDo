package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

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


func (s *Server) getTodosByDate(c *gin.Context) {
    rangeType := c.Query("range") // day/week/month
    dateStr := c.Query("date")    // YYYY-MM-DD format
    
    // Parse and validate date
    baseDate, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
        return
    }

    // Calculate date range with 1-year limit
    now := time.Now()
    if baseDate.Before(now.AddDate(-1, 0, 0)) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Date range exceeds 1 year limit"})
        return
    }

    var start, end time.Time
    switch rangeType {
    case "day":
        start = time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), 0, 0, 0, 0, time.UTC)
        end = start.AddDate(0, 0, 1)
    case "week":
        start = baseDate.AddDate(0, 0, -int(baseDate.Weekday()))
        end = start.AddDate(0, 0, 7)
    case "month":
        start = time.Date(baseDate.Year(), baseDate.Month(), 1, 0, 0, 0, 0, time.UTC)
        end = start.AddDate(0, 1, 0)
    default:
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid range type"})
        return
    }

    // Query todos within date range
    rows, err := s.db.Query(`
        SELECT id, title, description, completed, created_at, updated_at 
        FROM todos 
        WHERE created_at >= $1 AND created_at < $2
        ORDER BY created_at DESC
    `, start, end)
    
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

    // Group by day if weekly/monthly view
    if rangeType != "day" {
        grouped := make(map[string][]Todo)
        for _, todo := range todos {
            dateKey := todo.CreatedAt.Format("2006-01-02")
            grouped[dateKey] = append(grouped[dateKey], todo)
        }

        var result []GroupedTodos
        for date, items := range grouped {
            result = append(result, GroupedTodos{
                Date:  date,
                Todos: items,
            })
        }
        c.JSON(http.StatusOK, result)
        return
    }
    c.JSON(http.StatusOK, todos)
}
