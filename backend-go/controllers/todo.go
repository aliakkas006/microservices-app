package controllers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/aliakkas006/backend-go/db"
	"github.com/aliakkas006/backend-go/models"
	"github.com/gin-gonic/gin"
)

// GET /api/todos
func GetTodos(c *gin.Context) {
	rows, err := db.DB.Query(context.Background(), "SELECT id, title, description, completed FROM todos ORDER BY id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
		return
	}
	defer rows.Close()

	todos := []models.Todo{}
	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Completed); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse todos"})
			return
		}
		todos = append(todos, t)
	}

	c.JSON(http.StatusOK, todos)
}

// GET /api/todos/:id
func GetTodoByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var t models.Todo
	err = db.DB.QueryRow(context.Background(),
		"SELECT id, title, description, completed FROM todos WHERE id=$1", idInt).
		Scan(&t.ID, &t.Title, &t.Description, &t.Completed)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, t)
}

// POST /api/todos
func CreateTodo(c *gin.Context) {
	var t models.Todo
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if t.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	err := db.DB.QueryRow(context.Background(),
		"INSERT INTO todos (title, description, completed) VALUES ($1, $2, $3) RETURNING id",
		t.Title, t.Description, false).Scan(&t.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}
	t.Completed = false

	c.JSON(http.StatusCreated, t)
}

// PUT /api/todos/:id
func UpdateTodo(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var t models.Todo
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	commandTag, err := db.DB.Exec(context.Background(),
		"UPDATE todos SET title=$1, description=$2, completed=$3 WHERE id=$4",
		t.Title, t.Description, t.Completed, idInt)
	if err != nil || commandTag.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found or failed to update"})
		return
	}

	t.ID = idInt
	c.JSON(http.StatusOK, t)
}

// DELETE /api/todos/:id
func DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	commandTag, err := db.DB.Exec(context.Background(), "DELETE FROM todos WHERE id=$1", idInt)
	if err != nil || commandTag.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found or failed to delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}
