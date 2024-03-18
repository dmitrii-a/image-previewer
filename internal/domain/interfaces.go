package domain

// ImageRepository is an interface for image repository.
type ImageRepository interface {
	Save(image *Image) error
	Get(url string, width, height int) (*Image, error)
}

// ImageDownloader is an interface for image downloader.
type ImageDownloader interface {
	Download(url string, headers map[string][]string) (*Image, error)
}
