package service

import (
	"bytes"
	"image/jpeg"
	"math"
	"time"

	"snapcheck/internal/domain"
	"snapcheck/internal/exif"
	"snapcheck/internal/validator"
)

// AnalyzeImage coordinates the validation and analysis of a JPEG image.
func AnalyzeImage(data []byte, lastModified *time.Time) (*domain.AnalysisResult, error) {
	if err := validator.Validate(data); err != nil {
		return &domain.AnalysisResult{
			Error: err.Error(),
		}, nil
	}

	config, err := jpeg.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return &domain.AnalysisResult{
			Error: "failed to decode jpeg structure: " + err.Error(),
		}, nil
	}

	info := domain.ImageInfo{
		Width:     config.Width,
		Height:    config.Height,
		SizeBytes: int64(len(data)),
	}

	exifData, err := exif.Extract(data)
	if err != nil {
		exifData = &domain.ExifData{HasExif: false}
	}

	heuristics := domain.Heuristics{}

	if lastModified != nil && exifData.DateTimeOriginal != nil {
		diff := lastModified.Sub(*exifData.DateTimeOriginal)
		if math.Abs(diff.Seconds()) > 120 {
			heuristics.TimestampMismatch = true
		}
	}

	if exifData.Orientation == nil {
		heuristics.OrientationMissing = true
	}

	if info.Width < 200 || info.Height < 200 {
		heuristics.LowResolution = true
	}

	return &domain.AnalysisResult{
		ImageInfo:  info,
		Exif:       exifData,
		Heuristics: heuristics,
	}, nil
}
