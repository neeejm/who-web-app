package core

import "github.com/gofiber/fiber/v2"

func UseRoutes(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-type", "application/json; charset=utf-8")

		return c.JSON(struct {
			Msg string `json:"msg"`
		}{
			Msg: "Hello, World 👋!",
		})
	})

	app.Post("/drawbox", getImageBox)
}
