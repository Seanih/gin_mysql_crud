package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

type Task struct {
	TaskID    uint16 `json:"id"`
	TaskName  string `json:"task_name"`
	Completed bool   `json:"completed"`
	OwnerName string `json:"owner_name"`
}

func main() {
	dbConnection()

	app := gin.Default()
	port := "localhost:4000"

	app.GET("/", getAllTasks)
	app.POST("/", addTask)

	app.Run(port)

}

func dbConnection() {
	// load env variables
	err := godotenv.Load("../.env.local")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASSWORD"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "crud_data",
	}

	// connect to db
	db, err = sql.Open("mysql", cfg.FormatDSN())

	// handles errors; do something better for production
	if err != nil {
		log.Fatal(err)
	}

	// verify if db connection is live; reconnect if not
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("successfully connected to db!")
}

func getAllTasks(c *gin.Context) {
	rows, err := db.Query(("SELECT id, task_name, completed, owner_name FROM tasks JOIN owners ON tasks.owner_id = owners.owner_id"))

	if err != nil {
		log.Fatalf("error: %v", err)

	}
	defer rows.Close()

	var allTasks []Task

	for rows.Next() {
		var task Task

		if err := rows.Scan(&task.TaskID, &task.TaskName, &task.Completed, &task.OwnerName); err != nil {
			log.Fatalf("error: %v", err)
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
		log.Fatalf("error adding task: %v", err)
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
