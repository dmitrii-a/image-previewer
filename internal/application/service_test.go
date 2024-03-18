package application

import (
	"errors"
	"testing"

	"github.com/dmitrii-a/image-previewer/internal/infrastructure"
	"github.com/dmitrii-a/image-previewer/tests"
	"github.com/dmitrii-a/image-previewer/tests/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ImageServiceSuite struct {
	suite.Suite
	service    *ImageService
	repo       *mocks.ImageRepository
	downloader *mocks.ImageDownloader
}

func (s *ImageServiceSuite) SetupTest() {
	s.repo = &mocks.ImageRepository{}
	s.downloader = &mocks.ImageDownloader{}
	s.service = NewImageService(s.repo, s.downloader)
}

func (s *ImageServiceSuite) TestResizeImageSuccess() {
	s.repo.On("Get", "http://example.com/image.jpg", 100, 100).Return(tests.GenerateImage(), nil)
	data, err := s.service.ResizeImage("http://example.com/image.jpg", 100, 100, nil)
	s.Require().NoError(err)
	s.NotNil(data)
}

func (s *ImageServiceSuite) TestResizeImageDownloadError() {
	s.repo.On(
		"Get",
		"http://example.com/image.jpg", 100, 100,
	).Return(nil, infrastructure.NewErrImageNotFound(errors.New("error")))
	s.downloader.On(
		"Download",
		"http://example.com/image.jpg",
		map[string][]string{},
	).Return(nil, infrastructure.NewErrImageDownload(errors.New("error")))

	_, err := s.service.ResizeImage(
		"http://example.com/image.jpg", 100, 100, map[string][]string{},
	)

	s.Error(err)
}

func (s *ImageServiceSuite) TestResizeImageSaveError() {
	s.repo.On(
		"Get",
		"http://example.com/image.jpg",
		100,
		100,
	).Return(nil, infrastructure.NewErrImageNotFound(errors.New("error")))
	s.downloader.On(
		"Download",
		"http://example.com/image.jpg",
		map[string][]string{},
	).Return(tests.GenerateImage(), nil)
	s.repo.On(
		"Save", mock.Anything,
	).Return(infrastructure.NewErrImageSave(errors.New("error")))

	_, err := s.service.ResizeImage(
		"http://example.com/image.jpg",
		100,
		100, map[string][]string{},
	)

	s.Error(err)
}

func TestImageServiceSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ImageServiceSuite))
}
