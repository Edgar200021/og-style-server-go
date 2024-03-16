package services

import (
	"context"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type ImageUploaderService interface {
	Upload(file any) (string, error)
}

type CldImageUploaderService struct {
	Cloudinary *cloudinary.Cloudinary
}

func (c *CldImageUploaderService) Upload(file any) (string, error) {
	res, err := c.Cloudinary.Upload.Upload(context.Background(), file, uploader.UploadParams{
		Folder:         "og-style",
		UniqueFilename: api.Bool(true),
		UseFilename:    api.Bool(true),
	})

	if err != nil {
		return "", err
	}

	return res.SecureURL, nil
}
