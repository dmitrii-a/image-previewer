package infrastructure

import "github.com/dmitrii-a/image-previewer/internal/common"

// NewErrImageDownload returns a new instance of the image download error.
func NewErrImageDownload(err error) *common.ApplicationError {
	code := common.ImageDownloadErrorCode

	return &common.ApplicationError{
		Err:     err,
		Code:    code,
		Message: "image download error",
	}
}

// NewErrImageDecode returns a new instance of the image decode error.
func NewErrImageDecode(err error) *common.ApplicationError {
	code := common.ImageDecodeErrorCode

	return &common.ApplicationError{
		Err:     err,
		Code:    code,
		Message: "image decode error",
	}
}

// NewErrImageSave returns a new instance of the image save error.
func NewErrImageSave(err error) *common.ApplicationError {
	code := common.ImageSaveErrorCode

	return &common.ApplicationError{
		Err:     err,
		Code:    code,
		Message: "image download error",
	}
}

// NewErrImageGet returns a new instance of the image get error.
func NewErrImageGet(err error) *common.ApplicationError {
	code := common.ImageGetErrorCode

	return &common.ApplicationError{
		Err:     err,
		Code:    code,
		Message: "image get error",
	}
}

// NewErrImageNotFound returns a new instance of the image not found error.
func NewErrImageNotFound(err error) *common.ApplicationError {
	code := common.ImageNotFoundErrorCode

	return &common.ApplicationError{
		Err:     err,
		Code:    code,
		Message: "image not found",
	}
}
