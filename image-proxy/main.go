package main

import (
	"fmt"
	"proxy/core"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("__main__")

	// server
	app := fiber.New()
	core.UseRoutes(app)
	app.Listen(":3003")
}
