package domain

import (
	"image"

	"github.com/disintegration/imaging"
)

// Image struct for image.
type Image struct {
	Path string
	Data image.Image
	URL  string
}

// Resize resizes the image.
func (s *Image) Resize(width, height int) {
	s.Data = imaging.Resize(s.Data, width, height, imaging.Lanczos)
}
