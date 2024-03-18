package repository

import (
	"fmt"
	"image/jpeg"
	"os"
	"sync"

	"github.com/dmitrii-a/image-previewer/internal/common"
	"github.com/dmitrii-a/image-previewer/internal/domain"
	"github.com/dmitrii-a/image-previewer/internal/infrastructure"
	"github.com/google/uuid"
)

type imageCacheRepository struct {
	sync.RWMutex
	images      map[string]string
	keys        []string
	maxSize     int64
	currentSize int64
}

// NewImageCacheRepository returns a new instance of a imageCacheRepository.
func NewImageCacheRepository() domain.ImageRepository {
	return &imageCacheRepository{
		images:      make(map[string]string),
		keys:        []string{},
		maxSize:     common.Config.Cache.MaxSize,
		currentSize: 0,
	}
}

func (repo *imageCacheRepository) generateUUID() string {
	return uuid.New().String()
}

func (repo *imageCacheRepository) generateKey(url string, width, height int) string {
	return fmt.Sprintf("%s-%d-%d", url, width, height)
}

// Save saves the image to the cache.
func (repo *imageCacheRepository) Save(image *domain.Image) error {
	outFile, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("img-%v-*.jpg", repo.generateUUID()))
	if common.IsErr(err) {
		return infrastructure.NewErrImageSave(fmt.Errorf("error creating cached file: %w", err))
	}

	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {
			common.Logger.Error().Msgf("error closing file: %v", err)
		}
	}(outFile)

	if err := jpeg.Encode(outFile, image.Data, nil); err != nil {
		return infrastructure.NewErrImageSave(fmt.Errorf("error encoding image: %w", err))
	}

	stat, err := outFile.Stat()
	if common.IsErr(err) {
		return infrastructure.NewErrImageSave(fmt.Errorf("error getting file stat: %w", err))
	}

	repo.Lock()
	defer repo.Unlock()

	if repo.currentSize+stat.Size() > repo.maxSize && len(repo.keys) > 0 {
		key := repo.keys[0]

		err := os.Remove(repo.images[key])
		if common.IsErr(err) {
			return infrastructure.NewErrImageSave(fmt.Errorf("error removing file: %w", err))
		}

		delete(repo.images, repo.keys[0])
		repo.keys = repo.keys[1:]
	}

	key := repo.generateKey(image.URL, image.Data.Bounds().Dx(), image.Data.Bounds().Dy())

	repo.images[key] = outFile.Name()
	repo.keys = append(repo.keys, key)
	repo.currentSize += stat.Size()

	return nil
}

// Get returns the image from the cache.
func (repo *imageCacheRepository) Get(url string, width, height int) (*domain.Image, error) {
	repo.RLock()
	imgPath, ok := repo.images[repo.generateKey(url, width, height)]
	repo.RUnlock()

	if !ok {
		return nil, infrastructure.NewErrImageNotFound(fmt.Errorf("image not found: %s", url))
	}

	f, err := os.Open(imgPath)

	if common.IsErr(err) {
		return nil, infrastructure.NewErrImageGet(fmt.Errorf("error opening file: %w", err))
	}

	img, err := jpeg.Decode(f)
	if common.IsErr(err) {
		return nil, infrastructure.NewErrImageGet(fmt.Errorf("error decoding image: %w", err))
	}

	return &domain.Image{Data: img, URL: url}, nil
}
