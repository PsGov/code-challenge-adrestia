package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Global variable for db connection
var db *sqlx.DB

// User struct
type User struct {
	ID    int    `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
	Age   int    `json:"age" db:"age"`
}

func main() {
	// Initialize Fiber
	app := fiber.New()

	// Get environment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// Build connection string
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)

	// Connect to PostgreSQL
	var err error
	db, err = sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatal("Unable to connect to the database: ", err)
	}
	log.Println("Database connected")

	// Setup CORS
	app.Use(cors.New(cors.Config{
		// Allow Vue.js development port
		AllowOrigins: "http://localhost:9001",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	// Set up routes
	app.Get("/users", getUsers)
	app.Post("/users", createUser)
	app.Put("/users/:id", updateUser)
	app.Delete("/users/:id", deleteUser)

	// Start the server on port 9000
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "9000" // default value if not set
	}
	log.Fatal(app.Listen(":" + appPort))
}

// getUsers fetches users
func getUsers(c *fiber.Ctx) error {
	// Parse page and limit with defaults
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	// Get the search query from the request
	search := c.Query("search", "")

	// Calculate offset for pagination
	offset := (page - 1) * limit

	// Prepare the query with optional search filtering
	var users []User
	var query string
	var args []interface{}

	if search != "" {
		// Filter users by name or email when search term is present
		query = `SELECT * FROM users WHERE name ILIKE $1 OR email ILIKE $1 ORDER BY id LIMIT $2 OFFSET $3`
		args = append(args, "%"+search+"%", limit, offset)
	} else {
		// No search term; fetch all users for the page
		query = `SELECT * FROM users ORDER BY id LIMIT $1 OFFSET $2`
		args = append(args, limit, offset)
	}

	err = db.Select(&users, query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users"})
	}

	// Get the total count of users (with search filter applied)
	var totalCount int
	if search != "" {
		err = db.Get(&totalCount, `SELECT COUNT(*) FROM users WHERE name ILIKE $1 OR email ILIKE $1`, "%"+search+"%")
	} else {
		err = db.Get(&totalCount, `SELECT COUNT(*) FROM users`)
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count users"})
	}

	// Return users and total count
	return c.JSON(fiber.Map{
		"users":      users,
		"totalCount": totalCount,
	})
}

// createUser is used to create new users
func createUser(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		log.Println("Error parsing user data:", err)
		return c.Status(400).JSON(fiber.Map{"error": "Failed to parse request body"})
	}

	// Insert new user into the database
	_, err := db.Exec("INSERT INTO users (name, email, age) VALUES ($1, $2, $3)", user.Name, user.Email, user.Age)
	if err != nil {
		log.Println("Error creating user:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created successfully"})
}

// updateUser updates users based on id
func updateUser(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return err
	}
	id := c.Params("id")
	db.Exec("UPDATE users SET name=$1, email=$2, age=$3 WHERE id=$4", user.Name, user.Email, user.Age, id)
	return c.SendStatus(fiber.StatusOK)
}

// deleteUser deletes users based on id
func deleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	db.Exec("DELETE FROM users WHERE id=$1", id)
	return c.SendStatus(fiber.StatusOK)
}
