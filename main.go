package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"github.com/jinzhu/gorm"
)

func main() {

	router := gin.Default()

	v1 := router.Group("/paladin")
	{
		v1.POST("/", CreateTodo)
		v1.GET("/", FetchAllTodo)
		v1.GET("/:id", FetchSingleTodo)
		v1.PUT("/:id", UpdateTodo)
		v1.DELETE("/:id", DeleteTodo)
	}
	router.Run()

}

type Todo struct {
	gorm.Model
	Title     string `json:"title"`
	Completed int    `json:"completed"`
}

type TransformedTodo struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func CreateTodo(c *gin.Context) {
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	todo := Todo{Title: c.PostForm("title"), Completed: completed}
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!", "resourceId": todo.ID})
}

func FetchAllTodo(c *gin.Context) {
	var todos []Todo
	var _todos []TransformedTodo

	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	//transforms the todos for building a good response
	for _, item := range todos {
		completed := false
		if item.Completed == 1 {
			completed = true
		} else {
			completed = false
		}
		_todos = append(_todos, TransformedTodo{ID: item.ID, Title: item.Title, Completed: completed})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todos})
}

func FetchSingleTodo(c *gin.Context) {
	var todo Todo
	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	completed := false
	if todo.Completed == 1 {
		completed = true
	} else {
		completed = false
	}

	_todo := TransformedTodo{ID: todo.ID, Title: todo.Title, Completed: completed}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todo})
}

func UpdateTodo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo updated successfully!"})
}

func DeleteTodo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo deleted successfully!"})
}
