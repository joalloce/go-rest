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

	blogRepo := New()

	app.Post("/blogs", blogRepo.createBlog)
	app.Get("/blogs", blogRepo.getBlogs)
	app.Get("/blogs/:id", blogRepo.getBlog)
	app.Put("/blogs/:id", blogRepo.updateBlog)
	app.Delete("/blogs/:id", blogRepo.deleteBlogs)

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

type BlogRepo struct {
	Db *gorm.DB
}

func New() *BlogRepo {
	db := connectDB()
	db.AutoMigrate(&Blog{})
	return &BlogRepo{Db: db}
}

// create blog
func (repository *BlogRepo) createBlog(context *fiber.Ctx) error {
	var blog Blog
	err := context.BodyParser(&blog) // take teh data from body
	if err != nil {
		return context.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	err = repository.Db.Create(&blog).Error // create
	if err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return context.Status(fiber.StatusOK).JSON(blog)
}

// get blogs
func (repository *BlogRepo) getBlogs(context *fiber.Ctx) error {
	var blog []Blog
	err := repository.Db.Find(&blog).Error // get all blogs
	if err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return context.Status(fiber.StatusOK).JSON(blog)
}

// get blogs by id
func (repository *BlogRepo) getBlog(context *fiber.Ctx) error {
	var blog Blog
	id := context.Params("id")                                  // take id from params
	err := repository.Db.Where("id = ?", id).First(&blog).Error //find blog
	if err != nil {
		return context.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return context.Status(fiber.StatusOK).JSON(blog)
}

// update blog
func (repository *BlogRepo) updateBlog(context *fiber.Ctx) error {
	var blog Blog
	id := context.Params("id") // take id from params
	err := repository.Db.Where("id = ?", id).First(&blog).Error
	if err != nil {
		return context.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	err = context.BodyParser(&blog) // take the data from body
	if err != nil {
		return context.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	err = repository.Db.Save(&blog).Error // update user
	if err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return context.Status(fiber.StatusOK).JSON(blog)
}

// delete blog
func (repository *BlogRepo) deleteBlogs(context *fiber.Ctx) error {
	var blog Blog
	id := context.Params("id")                                   // take id from params
	err := repository.Db.Where("id = ?", id).Delete(&blog).Error // delete blog
	if err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return context.Status(fiber.StatusOK).JSON(fiber.Map{"message": "blog deleted"})
}
