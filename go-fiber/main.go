package main

import "github.com/gofiber/fiber/v2"

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
