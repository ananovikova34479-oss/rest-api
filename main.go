package main

import (
	"log"
	"net/http"
	"strconv" //  вместо "fmt" — для перевода числа в строку
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Task struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type TaskRequest struct {
	Title string `json:"title"`
}

var tasks = []Task{
	{ID: "1", Title: "Купить хлеб", Done: false},
	{ID: "2", Title: "Выучить Go", Done: true},
}

var nextID int = 3

func getTasks(c echo.Context) error {
	return c.JSON(http.StatusOK, tasks)
}

func createTask(c echo.Context) error {
	var req TaskRequest

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "неверный JSON"})
	}

	req.Title = strings.TrimSpace(req.Title)

	if req.Title == "" {

		return c.JSON(http.StatusBadRequest, map[string]string{"error": "пустой title"})
	}

	newTask := Task{ID: strconv.Itoa(nextID), Title: req.Title, Done: false}
	nextID++

	tasks = append(tasks, newTask)
	return c.JSON(http.StatusCreated, newTask)
}

func deleteTask(c echo.Context) error {
	id := c.Param("id")

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return c.NoContent(http.StatusNoContent)
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{"error": "задача не найдена"})
}

func main() {
	e := echo.New()

	e.Use(middleware.RequestLogger())

	e.GET("/tasks", getTasks)
	e.POST("/tasks", createTask)
	e.DELETE("/tasks/:id", deleteTask)

	err := e.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
