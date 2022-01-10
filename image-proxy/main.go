package main

import (
	"fmt"
	"proxy/core"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	fmt.Println("__main__")

	// server
	app := fiber.New()
	app.Use(recover.New())
	core.UseRoutes(app)
	app.Listen(":3003")
}
