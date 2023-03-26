package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Task struct {
	TaskID    uint16 `json:"id"`
	TaskName  string `json:"task_name"`
	Completed bool   `json:"completed"`
	OwnerName string `json:"owner_name"`
}

func getAllTasks(c *gin.Context) {
	rows, err := db.Query(("SELECT id, task_name, completed, owner_name FROM tasks JOIN owners ON tasks.owner_id = owners.owner_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	defer rows.Close()

	var allTasks []Task

	for rows.Next() {
		var task Task

		if err := rows.Scan(&task.TaskID, &task.TaskName, &task.Completed, &task.OwnerName); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}

		allTasks = append(allTasks, task)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusBadRequest, err)

	}

	c.JSON(http.StatusOK, allTasks)
}

func addTask(c *gin.Context) {
	var newTask Task

	if err := c.BindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	// must first add owner to owner table
	result1, err := db.Exec("INSERT INTO owners (owner_name) VALUES (?)", newTask.OwnerName)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	ownerID, err := result1.LastInsertId()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	// add ownerID as param for insert query
	result2, err := db.Exec("INSERT INTO tasks (task_name, completed, owner_id) VALUES (?, ?, ?)", newTask.TaskName, newTask.Completed, ownerID)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	newTaskID, err := result2.LastInsertId()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	response := fmt.Sprintf("task created with id: %v", newTaskID)

	c.JSON(http.StatusCreated, response)
}

func deleteTask(c *gin.Context) {
	var deletedTask Task
	if err := c.BindJSON(&deletedTask); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", deletedTask.TaskID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err})
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err})
	}

	c.JSON(http.StatusNoContent, nil)
}
