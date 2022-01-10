package utils

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

// get an image from the cloundinary cdn
// take the folder where the image is and the ID of the image
func GetImage(folderName string, publicID string) string {
	if folderName != "" {
		publicID = folderName + "/" + publicID
	}
	// get env variables
	env := GetENV()
	// Start by creating a new instance of Cloudinary using CLOUDINARY_URL environment variable.
	// Alternatively you can use cloudinary.NewFromParams() or cloudinary.NewFromURL().
	var cld, err = cloudinary.NewFromParams(env.CloudName, env.ApiKey, env.ApiSecret)
	if err != nil {
		log.Fatalf("Failed to intialize Cloudinary, %v", err)
	}

	var ctx = context.Background()

	asset, err := cld.Admin.Asset(ctx, admin.AssetParams{PublicID: publicID})
	if err != nil {
		log.Fatalf("Failed to get asset details, %v\n", err)
	}

	// Print some basic information about the asset.
	log.Printf("Public ID: %v, URL: %v\n", asset.PublicID, asset.SecureURL)
	return asset.SecureURL
}

// upload an image to the cloudinary cdn
// take the folder to upload to in the cdn, image(URL) to upload and the ID to use for the image in the cdn
func UploadImage(folderName string, imgURL string, publicID string) string {
	if folderName != "" {
		publicID = folderName + "/" + publicID
	}

	// get env variables
	env := GetENV()
	// Start by creating a new instance of Cloudinary using CLOUDINARY_URL environment variable.
	// Alternatively you can use cloudinary.NewFromParams() or cloudinary.NewFromURL().
	var cld, err = cloudinary.NewFromParams(env.CloudName, env.ApiKey, env.ApiSecret)
	if err != nil {
		log.Fatalf("Failed to intialize Cloudinary, %v", err)
	}

	var ctx = context.Background()

	// Upload an image to your Cloudinary account from a specified URL.
	uploadResult, err := cld.Upload.Upload(
		ctx,
		imgURL,
		uploader.UploadParams{PublicID: publicID})
	log.Println(publicID)
	if err != nil {
		log.Fatalf("Failed to upload file, %v\n", err)
	}

	log.Println(uploadResult.SecureURL)
	return uploadResult.SecureURL
}

func DownloadImage(URL, fileName string) error {
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
