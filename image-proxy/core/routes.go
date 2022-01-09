package core

import (
	ut "proxy/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	ib "github.com/neeejm/image-box/utils"
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

	return c.Status(fiber.StatusOK).JSON(struct {
		Msg string `json:"msg"`
	}{
		Msg: "Hello, World 👋!",
	})
}

func getImage(c *fiber.Ctx) error {
	URL := ut.GetImage(c.Params("folder"), c.Params("id"))

	return c.Status(fiber.StatusOK).JSON(struct {
		StatusCode int    `json:"status_code"`
		ImageURL   string `json:"image_url"`
	}{
		StatusCode: c.Response().StatusCode(),
		ImageURL:   URL,
	})
}

func uploadImage(c *fiber.Ctx) error {
	url := c.Query("url")
	folderName := c.Query("folder")

	if folderName == "" || url == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(struct {
			StatusCode  int    `json:"status_code"`
			Description string `json:"description"`
		}{
			StatusCode:  c.Response().StatusCode(),
			Description: "Bad reques. Not all query parameters are given.",
		})
	}

	id := uuid.New().String()
	imageURL := ut.UploadImage(folderName, url, id)

	return c.Status(fiber.StatusOK).JSON(struct {
		StatusCode int    `json:"status_code"`
		ImageURL   string `json:"image_url"`
	}{
		StatusCode: c.Response().StatusCode(),
		ImageURL:   imageURL,
	})
}

func drawBox(c *fiber.Ctx) error {
	url := c.Query("url")
	// check if url is given in the query param
	if url == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(struct {
			StatusCode  int    `json:"status_code"`
			Description string `json:"description"`
		}{
			StatusCode:  c.Response().StatusCode(),
			Description: "Bad reques. Not all query parameters are given.",
		})
	}

	// upload the image in the before folder
	id := uuid.New().String()
	imageBeforeURL := ut.UploadImage("before", url, id)

	// detect the face int the image
	box, err := ut.GetBoundingBox(ut.Fetch(imageBeforeURL))
	if err != nil {
		return c.Status(fiber.StatusAccepted).JSON(struct {
			StatusCode  int    `json:"status_code"`
			Description string `json:"description"`
		}{
			StatusCode:  c.Response().StatusCode(),
			Description: "Image not supported.",
		})
	}

	// download image
	ut.DownloadImage(imageBeforeURL, "face.png")

	// draw a box around the face
	ib.DrawBox("face.png", box)

	// upload the image in the after folder
	imageAfterURL := ut.UploadImage("after", "out.png", id)

	return c.Status(fiber.StatusOK).JSON(struct {
		StatusCode int    `json:"status_code"`
		ImageURL   string `json:"image_url"`
	}{
		StatusCode: c.Response().StatusCode(),
		ImageURL:   imageAfterURL,
	})
}
