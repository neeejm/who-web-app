package core

import (
	"fmt"
	"io/ioutil"
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

type Face struct {
	Image           []byte `json:"image"`
	Gender          bool   `json:"gender"`
	Name            string `json:"name"`
	ConfidenceLevel int    `json:"confidence_level"`
}
type Images struct {
	ImageBox     []byte   `json:"image_box"`
	Faces        []Face   `json:"faces"`
	BoundingBox  []ib.Box `json:"bounding_box"`
	CreationDate string   `json:"creation_date"`
}

type SuccessfulUpload struct {
	StatusCode  int    `json:"status_code"`
	Description string `json:"description"`
	Image       Image  `json:"image"`
}

type SuccessfulDrawing struct {
	StatusCode  int    `json:"status_code"`
	Description string `json:"description"`
	Images      Images `json:"images"`
}

type RequestError struct {
	StatusCode  int    `json:"status_code"`
	Description string `json:"description"`
}

var tmpDir string = "tmp/"

func UseRoutes(app *fiber.App) {
	// app.Use(logger.New())

	app.Get("/", HelloWorld)
	app.Get("/image/:folder/:id", getImage)
	app.Post("/upload", uploadImage)
	app.Post("/detect", detect)
	// app.Post("/find", find)
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

func detect(c *fiber.Ctx) error {
	defaultPNG := "face.png"
	defaultOutputPNG := "out.png"
	url := c.Query("url")

	// check if url is given in the query param
	if url == "" {
		return c.Status(fiber.ErrBadRequest.Code).JSON(RequestError{
			StatusCode:  c.Response().StatusCode(),
			Description: "Bad reques. Not all query parameters are given.",
		})
	}

	// detect the face int the image
	box, err := ut.GetBoundingBox(ut.Fetch(url, "face"))
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
	ut.DownloadImage(url, tmpDir+defaultPNG)

	// draw a box around the face
	ib.DrawBox(tmpDir+defaultPNG, box, tmpDir+defaultOutputPNG)
	imageBoxData, err := ioutil.ReadFile(tmpDir + defaultOutputPNG)
	if err != nil {
		return err
	}
	// return image buffer
	facesData := []Face{}
	// crop faces
	count := 0
	for i := 0; i < len(box); i++ {
		count++
		imgName := fmt.Sprintf(tmpDir+"cropped-face-%d.png", count)
		ib.CropImage(tmpDir+defaultPNG, box[i], imgName)
		data, err := ioutil.ReadFile(imgName)
		if err != nil {
			return err
		}
		id := uuid.New().String()
		imageURL := ut.UploadImage("faces", imgName, id)
		celebName, confidenceLevel, err := ut.GetCelebrity(ut.Fetch(imageURL, "celebrity"))
		gender, err := ut.GetGender(ut.Fetch(imageURL, "gender"))
		facesData = append(facesData, Face{
			Image:           data,
			Gender:          gender,
			Name:            celebName,
			ConfidenceLevel: confidenceLevel,
		})
		switch confidenceLevel {
		case 0:
			box[i].LineColor = "#ff0800"
		case 1:
			box[i].LineColor = "#fdbe02"
		case 2:
			box[i].LineColor = "#0cf20e"
		}
		// fmt.Println("byte slice data", data)
	}

	// id := uuid.New().String()
	// ut.UploadImage("after", tmpDir+defaultOutputPNG, id)
	// id = uuid.New().String()
	// ut.UploadImage("before", tmpDir+defaultPNG, id)

	return c.Status(fiber.StatusOK).JSON(SuccessfulDrawing{
		StatusCode:  c.Response().StatusCode(),
		Description: "Box drawen successfully.",
		Images: Images{
			ImageBox:     imageBoxData,
			Faces:        facesData,
			BoundingBox:  box,
			CreationDate: time.Now().Format("2006-01-02"),
		},
	})
}
