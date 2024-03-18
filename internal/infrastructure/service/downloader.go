package service

import (
	"context"
	"image"
	"net/http"

	"github.com/dmitrii-a/image-previewer/internal/common"
	"github.com/dmitrii-a/image-previewer/internal/domain"
	"github.com/dmitrii-a/image-previewer/internal/infrastructure"
)

type imageDownloaderService struct{}

// NewImageDownloaderService creates a new instance of the image downloader service.
func NewImageDownloaderService() domain.ImageDownloader {
	return &imageDownloaderService{}
}

// Download downloads the image from the given URL.
func (s *imageDownloaderService) Download(
	url string,
	headers map[string][]string,
) (*domain.Image, error) {
	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if common.IsErr(err) {
		return nil, err
	}

	for k, values := range headers {
		for _, value := range values {
			req.Header.Add(k, value)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if common.IsErr(err) {
		return nil, infrastructure.NewErrImageDownload(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, infrastructure.NewErrImageNotFound(err)
	}

	img, _, err := image.Decode(resp.Body)
	if common.IsErr(err) {
		return nil, infrastructure.NewErrImageDecode(err)
	}

	entity := &domain.Image{Data: img, URL: url}

	return entity, nil
}
