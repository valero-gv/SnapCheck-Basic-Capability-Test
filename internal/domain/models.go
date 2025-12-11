package domain

import "time"

// AnalysisResult represents the final JSON output structure.
type AnalysisResult struct {
	ImageInfo  ImageInfo  `json:"image_info"`
	Exif       *ExifData  `json:"exif,omitempty"`
	Heuristics Heuristics `json:"heuristics"`
	Error      string     `json:"error,omitempty"`
}

type ImageInfo struct {
	Width     int   `json:"width"`
	Height    int   `json:"height"`
	SizeBytes int64 `json:"size_bytes"`
}

type ExifData struct {
	HasExif          bool       `json:"has_exif"`
	DateTimeOriginal *time.Time `json:"date_time_original,omitempty"`
	Make             string     `json:"make,omitempty"`
	Model            string     `json:"model,omitempty"`
	Orientation      *int       `json:"orientation,omitempty"`
	GPS              *GPSData   `json:"gps,omitempty"`
}

type GPSData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Heuristics struct {
	TimestampMismatch  bool `json:"timestamp_mismatch"`
	OrientationMissing bool `json:"orientation_missing"`
	LowResolution      bool `json:"low_resolution"`
}

