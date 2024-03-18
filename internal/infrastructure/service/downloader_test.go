package service

import (
	"image"
	"image/jpeg"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmitrii-a/image-previewer/internal/domain"
	"github.com/stretchr/testify/suite"
)

type ImageDownloaderServiceTestSuite struct {
	suite.Suite
	service domain.ImageDownloader
	server  *httptest.Server
}

func (s *ImageDownloaderServiceTestSuite) SetupTest() {
	s.service = NewImageDownloaderService()
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		img := image.NewRGBA(image.Rect(0, 0, 100, 100))
		if err := jpeg.Encode(w, img, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(http.StatusOK)
	}))
}

func (s *ImageDownloaderServiceTestSuite) TearDownTest() {
	s.server.Close()
}

func (s *ImageDownloaderServiceTestSuite) TestSuccess() {
	img, err := s.service.Download(s.server.URL, make(map[string][]string))
	s.Require().NoError(err)
	s.NotNil(img)
}

func (s *ImageDownloaderServiceTestSuite) TestHTTPErrorCloseServer() {
	s.server.Close()
	img, err := s.service.Download(s.server.URL, make(map[string][]string))
	s.Require().Error(err)
	s.Nil(img)
}

func (s *ImageDownloaderServiceTestSuite) TestImageDecodeError() {
	s.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("not an image"))
	})
	img, err := s.service.Download(s.server.URL, make(map[string][]string))

	s.Require().Error(err)
	s.Nil(img)
}

func (s *ImageDownloaderServiceTestSuite) TestHTTPRequestError() {
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}))
	img, err := s.service.Download(s.server.URL, make(map[string][]string))
	s.Require().Error(err)
	s.Nil(img)
}

func TestImageDownloaderServiceTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ImageDownloaderServiceTestSuite))
}
