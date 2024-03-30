package main

import "github.com/gofiber/fiber/v2"

import (
    "myapp/code/utils"
    
)

func main() {
    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString(utils.TestUtilMessage())
    })

	app.Listen(":8080")
}