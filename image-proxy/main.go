package main

import (
	"fmt"

	ih "proxy/core"

	ib "github.com/neeejm/image-box/utils"
)

func main() {
	fmt.Println("__main__")

	box := ib.Box{
		TopRow:    0.17699452,
		RightCol:  0.6664675,
		BottomRow: 1,
		LeftCol:   0.35841435,
	}

	fileName := "face.png"
	URL := "https://www.thefaceparismanagement.com/uploads/images/products/1234.jpg"
	err := ih.DownloadFile(URL, fileName)
	if err != nil {
		panic(err)
	}

	ib.DrawBox("face.png", box)

	// app := fiber.New()
	// app.Use(logger.New())

	// useRoutes(app)

	// app.Listen(":3003")
}
