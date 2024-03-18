package tests

import (
	"image"

	"github.com/dmitrii-a/image-previewer/internal/domain"
)

func GenerateImage() *domain.Image {
	return &domain.Image{
		Data: image.NewRGBA(image.Rect(0, 0, 100, 100)),
	}
}
