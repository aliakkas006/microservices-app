package tests

import (
	"testing"

	"github.com/aliakkas006/todo-api/models"
)

func TestDummyTodoList(t *testing.T) {
	// dummy todos
	todos := []models.Todo{
		{ID: 1, Title: "Learn Go", Description: "Finish Go tutorial", Completed: false},
		{ID: 2, Title: "Build API", Description: "Create REST API in Go", Completed: true},
	}

	// check the length
	if len(todos) != 2 {
		t.Errorf("expected 2 todos, got %d", len(todos))
	}

	// check first todo is not completed
	if todos[0].Completed {
		t.Errorf("expected first todo to be not completed, got true")
	}

	// check second todo is completed
	if !todos[1].Completed {
		t.Errorf("expected second todo to be completed, got false")
	}
}
