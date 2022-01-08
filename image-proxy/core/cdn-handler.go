package core

import (
	"context"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/joho/godotenv"
)

type ENV struct {
	cloudName string
	apiKey    string
	apiSecret string
}

func getENV() ENV {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env := ENV{
		cloudName: os.Getenv("CLOUD_NAME"),
		apiKey:    os.Getenv("API_KEY"),
		apiSecret: os.Getenv("API_SECRET"),
	}

	return env
}

func GetImage(folderName string, publicID string) {
	if folderName != "" {
		publicID = folderName + "/" + publicID
	}
	// get env variables
	env := getENV()
	// Start by creating a new instance of Cloudinary using CLOUDINARY_URL environment variable.
	// Alternatively you can use cloudinary.NewFromParams() or cloudinary.NewFromURL().
	var cld, err = cloudinary.NewFromParams(env.cloudName, env.apiKey, env.apiSecret)
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
}

func UploadImage(folderName string, imgURL string, publicID string) {
	if folderName != "" {
		publicID = folderName + "/" + publicID
	}

	// get env variables
	env := getENV()
	// Start by creating a new instance of Cloudinary using CLOUDINARY_URL environment variable.
	// Alternatively you can use cloudinary.NewFromParams() or cloudinary.NewFromURL().
	var cld, err = cloudinary.NewFromParams(env.cloudName, env.apiKey, env.apiSecret)
	if err != nil {
		log.Fatalf("Failed to intialize Cloudinary, %v", err)
	}

	var ctx = context.Background()

	// Upload an image to your Cloudinary account from a specified URL.
	uploadResult, err := cld.Upload.Upload(
		ctx,
		imgURL,
		uploader.UploadParams{PublicID: publicID})
	if err != nil {
		log.Fatalf("Failed to upload file, %v\n", err)
	}

	log.Println(uploadResult.SecureURL)
}
