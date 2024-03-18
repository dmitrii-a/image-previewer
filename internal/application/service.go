package application

import (
	"bytes"
	"image/jpeg"

	"github.com/dmitrii-a/image-previewer/internal/common"
	"github.com/dmitrii-a/image-previewer/internal/domain"
)

// ImageService struct for image service.
type ImageService struct {
	repository             domain.ImageRepository
	imageDownloaderService domain.ImageDownloader
}

// NewImageService returns a new instance of the image service.
func NewImageService(
	repository domain.ImageRepository,
	imageDownloaderService domain.ImageDownloader,
) *ImageService {
	return &ImageService{repository: repository, imageDownloaderService: imageDownloaderService}
}

// ResizeImage resizes the image.
func (s *ImageService) ResizeImage(
	url string,
	width, height int,
	headers map[string][]string,
) ([]byte, error) {
	var (
		img *domain.Image
		err error
	)

	img, err = s.repository.Get(url, width, height)
	if common.IsErr(err) {
		common.Logger.Debug().Msgf("error getting image from cache: %v", err)

		img, err = s.imageDownloaderService.Download(url, headers)
		if common.IsErr(err) {
			return nil, err
		}

		img.Resize(width, height)

		err = s.repository.Save(img)
		if common.IsErr(err) {
			common.Logger.Error().Msgf("error saving image to cache: %v", err)

			return nil, err
		}
	}

	buf := new(bytes.Buffer)

	err = jpeg.Encode(buf, img.Data, nil)
	if common.IsErr(err) {
		common.Logger.Error().Msgf("error encoding image: %v", err)

		return nil, err
	}

	return buf.Bytes(), nil
}
