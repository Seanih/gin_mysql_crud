package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {
	dbConnection()

	app := fiber.New()
	port := ":4000"

	fiberTest := fmt.Sprintf("server running on port%v", port)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiberTest)
	})

	log.Fatal(app.Listen(port))
}

func dbConnection() {
	// load env variables
	err := godotenv.Load()
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

	fmt.Println("Connected!")
}
