package helper

import (
	"context"
	"strings"
	"time"

	"e-commerce-api/config"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadFile(file interface{}) (string, error) {
	cld, err := cloudinary.NewFromParams(config.CloudinaryName, config.CloudinaryApiKey, config.CloudinaryApiScret)
	if err != nil {
		return "", err
	}

	publicID := time.Now().Format("20060102-150405") // Format  "(YY-MM-DD)-(hh-mm-ss)""

	uploadResult, err := cld.Upload.Upload(
		context.Background(),
		file,
		uploader.UploadParams{
			PublicID:     publicID,
			ResourceType: "image",
			Folder:       config.CloudinaryUploadFolder,
		})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}

func GetPublicID(secureURL string) string {
	// Proses filter Public ID dari SecureURL(avatar)
	urls := strings.Split(secureURL, "/")
	urls = urls[len(urls)-3:]                               // array [file, user, random_name.extension]
	noExtension := strings.Split(urls[len(urls)-1], ".")[0] // remove ".extension", result "random_name"
	urls = append(urls[:2], noExtension)                    // new array [file, user, random_name]
	publicID := strings.Join(urls, "/")                     // "file/user/random_name"

	return publicID
}

func DestroyFile(publicID string) error {
	cld, err := cloudinary.NewFromParams(config.CloudinaryName, config.CloudinaryApiKey, config.CloudinaryApiScret)
	if err != nil {
		return err
	}

	// Proses destroy file
	_, err = cld.Upload.Destroy(
		context.Background(),
		uploader.DestroyParams{
			PublicID: publicID,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
