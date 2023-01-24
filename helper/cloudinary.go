package helper

import (
	"context"
	"errors"
	"strings"
	"time"

	"e-commerce-api/config"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadFile(file interface{}) (string, error) {
	cld, err := cloudinary.NewFromParams(config.CLOUDINARY_CLOUD_NAME, config.CLOUDINARY_API_KEY, config.CLOUDINARY_API_SECRET)
	if err != nil {
		return "", err
	}
	// cld, err := cloudinary.NewFromParams("dnji8pgyl", "455139436831234", "iMQL2OVjJQKttb05He-r7pBN_8k")
	// if err != nil {
	// 	return "", err
	// }

	publicID := time.Now().Format("20060102-150405") // Format  "(YY-MM-DD)-(hh-mm-ss)""

	uploadResult, err := cld.Upload.Upload(
		context.Background(),
		file,
		uploader.UploadParams{
			PublicID:     publicID,
			ResourceType: "image",
			Folder:       config.CLOUDINARY_UPLOAD_FOLDER,
		})
	if err != nil {
		return "", err
	}

	// Pengecekan format file
	format := uploadResult.Format
	if format != "jpeg" && format != "png" && format != "jpg" {
		cld.Upload.Destroy(
			context.Background(),
			uploader.DestroyParams{
				PublicID: uploadResult.PublicID,
			},
		)
		return "", errors.New("kesalahan input user karena format gambar bukan jpg, jpeg, ataupun png")
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
	cld, err := cloudinary.NewFromParams(config.CLOUDINARY_CLOUD_NAME, config.CLOUDINARY_API_KEY, config.CLOUDINARY_API_SECRET)
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
