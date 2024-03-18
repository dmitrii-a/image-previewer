package common

import (
	"fmt"
)

const (
	ImageDownloadErrorCode = iota
	ImageDecodeErrorCode
	ImageSaveErrorCode
	ImageGetErrorCode
	ImageNotFoundErrorCode
)

// ApplicationError Typed Error with necessary fields.
type ApplicationError struct {
	Err     error
	Code    int
	Message string
}

// Error returns the error message.
func (e *ApplicationError) Error() string {
	return fmt.Sprintf("%s (%s)", e.Message, e.Err)
}
