package routes

import (
	"github.com/aliakkas006/todo-api/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterTodoRoutes(router *gin.Engine) {
	todoRoutes := router.Group("/api/todos")
	{
		todoRoutes.GET("", controllers.GetTodos)
		todoRoutes.GET("/:id", controllers.GetTodoByID)
		todoRoutes.POST("", controllers.CreateTodo)
		todoRoutes.PUT("/:id", controllers.UpdateTodo)
		todoRoutes.DELETE("/:id", controllers.DeleteTodo)
	}
}
