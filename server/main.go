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

	app.GET("/", func(ctx *gin.Context) {
		tasks, err := getAllTasks()

		if err != nil {
			log.Fatal(err)
		}
		ctx.IndentedJSON(http.StatusOK, tasks)
	})

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

func getAllTasks() ([]Task, error) {
	rows, err := db.Query(("SELECT id, task_name, completed, owner_name FROM tasks JOIN owners ON tasks.owner_id = owners.owner_id"))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allTasks []Task

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.TaskID, &task.TaskName, &task.Completed, &task.OwnerName); err != nil {
			return allTasks, err
		}

		allTasks = append(allTasks, task)
	}

	if err = rows.Err(); err != nil {
		return allTasks, err
	}

	return allTasks, nil
}

func addTask(newTask Task) (int64, error) {
	result, err := db.Exec("INSERT INTO owners (owner_name) VALUES (?)", newTask.OwnerName)

	if err != nil {
		return 0, fmt.Errorf("there was an error: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("there was an error: %v", err)
	}

	return id, nil
}
