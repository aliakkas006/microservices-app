package tests

import (
	"testing"

	"github.com/aliakkas006/email-service/models"
)

func TestDummyTodoList(t *testing.T) {

	todos := []models.Todo{
		{ID: 1, Title: "Learn Go", Description: "Finish Go tutorial", Completed: false, UserEmail: "ali@example.com"},
		{ID: 2, Title: "Build API", Description: "Create REST API", Completed: true, UserEmail: "john.doe@example.com"},
		{ID: 3, Title: "Write Tests", Description: "Write basic tests", Completed: false, UserEmail: "jane@example.com"},
	}

	// check total todos
	if len(todos) != 3 {
		t.Errorf("expected 3 todos, got %d", len(todos))
	}

	// check first todo's Completed and UserEmail
	if todos[0].Completed {
		t.Errorf("expected first todo to be not completed")
	}
	if todos[0].UserEmail != "ali@example.com" {
		t.Errorf("expected first todo UserEmail to be 'ali@example.com', got '%s'", todos[0].UserEmail)
	}

	// check second todo is completed
	if !todos[1].Completed {
		t.Errorf("expected second todo to be completed")
	}

	// check third todo Title
	if todos[2].Title != "Write Tests" {
		t.Errorf("expected third todo Title to be 'Write Tests', got '%s'", todos[2].Title)
	}
}
