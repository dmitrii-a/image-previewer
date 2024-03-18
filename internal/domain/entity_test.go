package domain

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func generateImage() *Image {
	return &Image{
		Data: image.NewRGBA(image.Rect(0, 0, 100, 100)),
	}
}

func TestImageResize(t *testing.T) {
	t.Parallel()

	img := generateImage()
	img.Resize(50, 50)

	assert.Equal(t, 50, img.Data.Bounds().Dx())
	assert.Equal(t, 50, img.Data.Bounds().Dy())
}

func TestImageResizeZeroWidthHeight(t *testing.T) {
	t.Parallel()

	img := generateImage()
	img.Resize(0, 0)

	assert.Equal(t, 0, img.Data.Bounds().Dx())
	assert.Equal(t, 0, img.Data.Bounds().Dy())
}

func TestImageResizeNegativeWidthHeight(t *testing.T) {
	t.Parallel()

	img := generateImage()
	img.Resize(-100, -100)

	assert.Equal(t, 0, img.Data.Bounds().Dx())
	assert.Equal(t, 0, img.Data.Bounds().Dy())
}

func TestImageIncreaseSize(t *testing.T) {
	t.Parallel()

	img := generateImage()
	img.Resize(300, 300)

	assert.Equal(t, 300, img.Data.Bounds().Dx())
	assert.Equal(t, 300, img.Data.Bounds().Dy())
}
