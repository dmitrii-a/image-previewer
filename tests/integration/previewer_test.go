package integration_test

import (
	"bufio"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"testing"

	"github.com/dmitrii-a/image-previewer/internal/common"
	. "github.com/onsi/ginkgo" //nolint: revive
	. "github.com/onsi/gomega" //nolint: revive
)

const ImageURL = "https://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg" //nolint: lll

var _ = Describe("Image Previewer", func() {
	serverURL := fmt.Sprintf(
		"http://%v", common.GetServerAddr(common.Config.Server.Host, common.Config.Server.Port),
	)

	Describe("Resize image", func() {
		Context("when the image exists", func() {
			It("should resize the image successfully", func() {
				resp, err := http.Get(fmt.Sprintf("%s/fill/300/200/%s", serverURL, ImageURL))
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				img, err := jpeg.Decode(bufio.NewReader(resp.Body))
				Expect(err).NotTo(HaveOccurred())

				bounds := img.Bounds()
				Expect(bounds.Dx()).To(Equal(300))
				Expect(bounds.Dy()).To(Equal(200))
			})
		})

		Context("when the image does not exist", func() {
			It("should return an error", func() {
				resp, err := http.Get(
					fmt.Sprintf(
						"%s/fill/300/200/https://raw.githubusercontent.com/OtusGolang/"+
							"final_project/master/examples/image-previewer/"+
							"nonexistent.jpg", serverURL,
					),
				)
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
						common.Logger.Error().Msgf("failed to close response body: %v", err)
					}
				}(resp.Body)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})

		Context("when the remote server returns an error", func() {
			It("should return an error", func() {
				resp, err := http.Get(
					fmt.Sprintf("%s/fill/300/200/http://127.0.0.1/image.jpg", serverURL),
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusBadGateway))
			})
		})

		Context("when the image is smaller than the required size", func() {
			It("should return correct result", func() {
				resp, err := http.Get(fmt.Sprintf("%s/fill/3000/2000/%s", serverURL, ImageURL))
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				img, err := jpeg.Decode(bufio.NewReader(resp.Body))
				Expect(err).NotTo(HaveOccurred())

				bounds := img.Bounds()
				Expect(bounds.Dx()).To(Equal(3000))
				Expect(bounds.Dy()).To(Equal(2000))
			})
		})

		Context("when width and height are incorrect", func() {
			It("should return an error for width param", func() {
				resp, err := http.Get(fmt.Sprintf("%s/fill/a/0/%s", serverURL, ImageURL))
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})
			It("should return an error for height param", func() {
				resp, err := http.Get(fmt.Sprintf("%s/fill/100/a/%s", serverURL, ImageURL))
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})
	})
})

func TestImageService(t *testing.T) {
	t.Parallel()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Image Previewer Integration Tests")
}
