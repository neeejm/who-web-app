package main

import (
	"fmt"

	ih "proxy/core"
	// ib "github.com/neeejm/image-box/utils"
)

func main() {
	fmt.Println("__main__")

	// cdn
	// ih.UploadImage("before", "https://i0.wp.com/post.medicalnewstoday.com/wp-content/uploads/sites/3/2020/03/GettyImages-1092658864_hero-1024x575.jpg?w=1155&h=1528", "face2")
	ih.GetImage("before", "face2")

	// image-box
	// box := ib.Box{
	// 	TopRow:    0.17699452,
	// 	RightCol:  0.6664675,
	// 	BottomRow: 1,
	// 	LeftCol:   0.35841435,
	// }

	// fileName := "face.png"
	// URL := "https://www.thefaceparismanagement.com/uploads/images/products/1234.jpg"
	// err := ih.DownloadFile(URL, fileName)
	// if err != nil {
	// 	panic(err)
	// }

	// ib.DrawBox("face.png", box)

	// server
	// app := fiber.New()
	// app.Use(logger.New())

	// useRoutes(app)

	// app.Listen(":3003")
}
