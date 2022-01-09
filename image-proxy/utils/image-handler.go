package utils

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

type BoundingBox struct {
	TopRow    float64 `json:"top_row"`
	RightCol  float64 `json:"right_col"`
	BottomRow float64 `json:"bottom_row"`
	LeftCol   float64 `json:"left_col"`
}

type Payload struct {
	BoundingBox BoundingBox
	ImageUrl    string `json:"image_url"`
}

func GetImageBox(c *fiber.Ctx) error {
	payload := new(Payload)
	err := c.BodyParser(payload)
	if err != nil {
		c.Status(fiber.StatusBadRequest).SendString(err.Error())
		return err
	}

	// box := ib.Box{
	// 	TopRow:    payload.BoundingBox.TopRow,
	// 	RightCol:  payload.BoundingBox.RightCol,
	// 	BottomRow: payload.BoundingBox.BottomRow,
	// 	LeftCol:   payload.BoundingBox.LeftCol,
	// }

	return c.Status(fiber.StatusOK).JSON(payload)
}

func DownloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
