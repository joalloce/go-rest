package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	app := setupRouter()
	app.Listen(":3000")
}

func setupRouter() *fiber.App {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	return app
}

const DB_USERNAME = "root"
const DB_PASSWORD = "password"
const DB_NAME = "db"
const DB_HOST = "db"
const DB_PORT = "3306"

// db connection
func connectDB() *gorm.DB {
	var err error
	DB_URI := DB_USERNAME + ":" + DB_PASSWORD + "@tcp" + "(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?" + "parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(DB_URI), &gorm.Config{})

	if err != nil {
		fmt.Println("Error db :", err)
		return nil
	}

	return db
}

// Blog model
type Blog struct {
	gorm.Model
	ID      int
	Author  int
	Upvotes int
}

type UserRepo struct {
	Db *gorm.DB
}

func New() *UserRepo {
	db := connectDB()
	db.AutoMigrate(&Blog{})
	return &UserRepo{Db: db}
}
