package core

import (
	ut "proxy/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UseRoutes(app *fiber.App) {
	// app.Use(logger.New())

	app.Get("/", HelloWorld)
	app.Get("/image/:folder/:id", getImage)
	app.Post("/upload", uploadImage)
	app.Post("/drawbox", drawBox)
}

func HelloWorld(c *fiber.Ctx) error {
	c.Set("Content-type", "application/json; charset=utf-8")

	return c.JSON(struct {
		Msg string `json:"msg"`
	}{
		Msg: "Hello, World 👋!",
	})
}

func getImage(c *fiber.Ctx) error {
	URL := ut.GetImage(c.Params("folder"), c.Params("id"))

	return c.JSON(struct {
		StatusCode int    `json:"status_code"`
		ImageURL   string `json:"image_url"`
	}{
		StatusCode: c.Response().StatusCode(),
		ImageURL:   URL,
	})
}

func uploadImage(c *fiber.Ctx) error {
	if c.Query("folder") == "" && c.Query("id") == "" {
		return c.SendStatus(400)
		// return c.JSON(struct {
		// 	StatusCode  int    `json:"status_code"`
		// 	Description string `json:"description"`
		// }{
		// 	StatusCode:  c.Status(400),
		// 	Description: "No request parameters",
		// })

	}

	id := uuid.New().String()
	URL := ut.UploadImage(c.Params("folder"), c.Params("url"), id)

	return c.JSON(struct {
		StatusCode int    `json:"status_code"`
		ImageURL   string `json:"image_url"`
	}{
		StatusCode: c.Response().StatusCode(),
		ImageURL:   URL,
	})
}

func drawBox(c *fiber.Ctx) error {

	return c.JSON(struct {
		Msg string `json:"msg"`
	}{
		Msg: "",
	})
}
