package validator

import (
	"bytes"
	"errors"
	"net/http"
)

var (
	ErrInvalidSize = errors.New("file size must be between 2KB and 20MB")
	ErrNotJPEG     = errors.New("file is not a valid JPEG")
	ErrCorrupted   = errors.New("file is corrupted or truncated")
	ErrEmpty       = errors.New("file is empty")
	ErrInvalidMIME = errors.New("invalid mime type, expected image/jpeg")
)

const (
	MinSize = 2 * 1024
	MaxSize = 20 * 1024 * 1024
)

// Validate performs binary-level validation on the image data.
func Validate(data []byte) error {
	if len(data) == 0 {
		return ErrEmpty
	}

	if len(data) < MinSize || len(data) > MaxSize {
		return ErrInvalidSize
	}

	if !bytes.HasPrefix(data, []byte{0xFF, 0xD8}) {
		return ErrNotJPEG
	}

	if !bytes.HasSuffix(data, []byte{0xFF, 0xD9}) {
		return ErrCorrupted
	}

	mimeType := http.DetectContentType(data)
	if mimeType != "image/jpeg" {
		return ErrInvalidMIME
	}

	return nil
}
