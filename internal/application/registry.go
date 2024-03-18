package application

import (
	"github.com/dmitrii-a/image-previewer/internal/infrastructure/repository"
	infraService "github.com/dmitrii-a/image-previewer/internal/infrastructure/service"
)

// ImageApplicationService is an instance of the image service.
var ImageApplicationService *ImageService

func init() {
	imageRepository := repository.NewImageCacheRepository()
	imageDownloaderService := infraService.NewImageDownloaderService()
	ImageApplicationService = NewImageService(imageRepository, imageDownloaderService)
}
