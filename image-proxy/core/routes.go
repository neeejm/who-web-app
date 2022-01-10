package core

import (
	ut "proxy/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	ib "github.com/neeejm/image-box/utils"
)

type Image struct {
	UUID         string `json:"uuid"`
	URL          string `json:"url"`
	Folder       string `json:"folder"`
	CreationDate string `json:"creation_date"`
}

type ImageBox struct {
	URL          string `json:"url"`
	BoundingBox  ib.Box `json:"bounding_box"`
	CreationDate string `json:"creation_date"`
}

type SuccessfulUpload struct {
	StatusCode  int    `json:"status_code"`
	Description string `json:"description"`
	Image       Image  `json:"image"`
}

type SuccessfulDrawing struct {
	StatusCode  int      `json:"status_code"`
	Description string   `json:"description"`
	ImageBox    ImageBox `json:"image"`
}

type RequestError struct {
	StatusCode  int    `json:"status_code"`
	Description string `json:"description"`
}

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
	if URL == "" {
		return c.Status(fiber.StatusAccepted).JSON(struct {
			StatusCode  int    `json:"status_code"`
			Description string `json:"description"`
		}{
			StatusCode:  c.Response().StatusCode(),
			Description: "No image found.",
		})
	}

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
		return c.Status(fiber.ErrBadRequest.Code).JSON(RequestError{
			StatusCode:  c.Response().StatusCode(),
			Description: "Bad reques. Not all query parameters are given.",
		})
	}

	id := uuid.New().String()
	imageURL := ut.UploadImage(folderName, url, id)

	return c.Status(fiber.StatusOK).JSON(SuccessfulUpload{
		StatusCode:  c.Response().StatusCode(),
		Description: "Image uploades successfully.",
		Image: Image{
			UUID:         id,
			URL:          imageURL,
			Folder:       folderName,
			CreationDate: time.Now().Format("2006-01-02"),
		},
	})
}

func drawBox(c *fiber.Ctx) error {
	url := c.Query("url")

	// check if url is given in the query param
	if url == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(RequestError{
			StatusCode:  c.Response().StatusCode(),
			Description: "Bad reques. Not all query parameters are given.",
		})
	}

	// detect the face int the image
	box, err := ut.GetBoundingBox(ut.Fetch(url))
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
	ut.DownloadImage(url, "face.png")

	// draw a box around the face
	ib.DrawBox("face.png", box)

	return c.Status(fiber.StatusOK).JSON(SuccessfulDrawing{
		StatusCode:  c.Response().StatusCode(),
		Description: "Box drawen successfully.",
		ImageBox: ImageBox{
			URL:          "",
			BoundingBox:  box,
			CreationDate: time.Now().Format("2006-01-02"),
		},
	})
}
