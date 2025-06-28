package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strconv"
    "github.com/thakurnishu/MinimalDo/database"
    "github.com/thakurnishu/MinimalDo/models"

    "github.com/gorilla/mux"
)

func GetTodos(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    rows, err := database.DB.Query(`
        SELECT id, title, description, completed, created_at, updated_at 
        FROM todos 
        ORDER BY created_at DESC
    `)
    if err != nil {
        http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var todos []models.Todo
    for rows.Next() {
        var todo models.Todo
        err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, 
            &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
        if err != nil {
            http.Error(w, "Failed to scan todo", http.StatusInternalServerError)
            return
        }
        todos = append(todos, todo)
    }

    json.NewEncoder(w).Encode(todos)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    var req models.CreateTodoRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if req.Title == "" {
        http.Error(w, "Title is required", http.StatusBadRequest)
        return
    }

    var todo models.Todo
    err := database.DB.QueryRow(`
        INSERT INTO todos (title, description) 
        VALUES ($1, $2) 
        RETURNING id, title, description, completed, created_at, updated_at
    `, req.Title, req.Description).Scan(
        &todo.ID, &todo.Title, &todo.Description, 
        &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

    if err != nil {
        http.Error(w, "Failed to create todo", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(todo)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid todo ID", http.StatusBadRequest)
        return
    }

    var req models.UpdateTodoRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Build dynamic update query
    query := "UPDATE todos SET updated_at = CURRENT_TIMESTAMP"
    args := []interface{}{}
    argCount := 0

    if req.Title != nil {
        argCount++
        query += ", title = $" + strconv.Itoa(argCount)
        args = append(args, *req.Title)
    }

    if req.Description != nil {
        argCount++
        query += ", description = $" + strconv.Itoa(argCount)
        args = append(args, *req.Description)
    }

    if req.Completed != nil {
        argCount++
        query += ", completed = $" + strconv.Itoa(argCount)
        args = append(args, *req.Completed)
    }

    argCount++
    query += " WHERE id = $" + strconv.Itoa(argCount) + " RETURNING id, title, description, completed, created_at, updated_at"
    args = append(args, id)

    var todo models.Todo
    err = database.DB.QueryRow(query, args...).Scan(
        &todo.ID, &todo.Title, &todo.Description, 
        &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Todo not found", http.StatusNotFound)
            return
        }
        http.Error(w, "Failed to update todo", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(todo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid todo ID", http.StatusBadRequest)
        return
    }

    result, err := database.DB.Exec("DELETE FROM todos WHERE id = $1", id)
    if err != nil {
        http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        http.Error(w, "Todo not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
