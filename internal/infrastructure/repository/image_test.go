package repository

import (
	"image"
	"os"
	"testing"

	"github.com/dmitrii-a/image-previewer/internal/common"
	"github.com/dmitrii-a/image-previewer/internal/domain"
	"github.com/stretchr/testify/suite"
)

type ImageCacheRepositorySuite struct {
	suite.Suite
	repo *imageCacheRepository
	img  *domain.Image
}

func (s *ImageCacheRepositorySuite) SetupTest() {
	s.repo = &imageCacheRepository{
		images:      make(map[string]string),
		keys:        []string{},
		maxSize:     common.Config.Cache.MaxSize,
		currentSize: 0,
	}
	s.img = &domain.Image{
		Data: image.NewRGBA(image.Rect(0, 0, 100, 100)),
		URL:  "http://example.com/image.jpg",
	}
}

func (s *ImageCacheRepositorySuite) TestGetImageExists() {
	err := s.repo.Save(s.img)
	s.Require().NoError(err)

	result, err := s.repo.Get(s.img.URL, s.img.Data.Bounds().Dx(), s.img.Data.Bounds().Dy())
	s.Require().NoError(err)
	s.Equal(s.img.URL, result.URL)
}

func (s *ImageCacheRepositorySuite) TestGetImageDoesNotExist() {
	_, err := s.repo.Get("http://example.com/image.jpg", 100, 100)
	s.Error(err)
}

func (s *ImageCacheRepositorySuite) TestGetDeletedImage() {
	err := s.repo.Save(s.img)
	s.Require().NoError(err)

	err = os.Remove(s.repo.keys[0])
	if err != nil {
		return
	}

	_, err = s.repo.Get(s.img.URL, s.img.Data.Bounds().Dx(), s.img.Data.Bounds().Dy())
	s.Require().Error(err)
}

func (s *ImageCacheRepositorySuite) TestSaveImage() {
	err := s.repo.Save(s.img)
	s.Require().NoError(err)

	result, err := s.repo.Get(s.img.URL, s.img.Data.Bounds().Dx(), s.img.Data.Bounds().Dy())
	s.Require().NoError(err)
	s.Equal(s.img.URL, result.URL)
}

func (s *ImageCacheRepositorySuite) TestSaveImageExceedsMaxSize() {
	s.repo.maxSize = 1
	img2 := &domain.Image{
		Data: image.NewRGBA(image.Rect(0, 0, 100, 100)),
		URL:  "http://example.com/image2.jpg",
	}

	err := s.repo.Save(s.img)
	s.Require().NoError(err)

	err = s.repo.Save(img2)
	s.Require().NoError(err)

	_, err = s.repo.Get(s.img.URL, s.img.Data.Bounds().Dx(), s.img.Data.Bounds().Dy())
	s.Require().Error(err)
}

func (s *ImageCacheRepositorySuite) TestSaveImageNonExistFolder() {
	err := os.Setenv("TMPDIR", "/nonexistent")
	defer os.Unsetenv("TMPDIR")
	s.Require().NoError(err)

	err = s.repo.Save(s.img)
	s.Require().Error(err)
}

func TestImageCacheRepositorySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ImageCacheRepositorySuite))
}
