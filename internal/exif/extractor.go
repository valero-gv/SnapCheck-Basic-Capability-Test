package exif

import (
	"time"

	"snapcheck/internal/domain"

	exif "github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
)

// Extract parses the image data and returns EXIF information.
// It handles missing or corrupted EXIF data gracefully.
func Extract(data []byte) (*domain.ExifData, error) {
	rawExif, err := exif.SearchAndExtractExif(data)
	if err != nil {
		return &domain.ExifData{HasExif: false}, nil
	}

	im, err := exifcommon.NewIfdMappingWithStandard()
	if err != nil {
		return &domain.ExifData{HasExif: false}, nil
	}

	ti := exif.NewTagIndex()
	_, index, err := exif.Collect(im, ti, rawExif)
	if err != nil {
		return &domain.ExifData{HasExif: false}, nil
	}

	out := &domain.ExifData{HasExif: true}
	rootIfd := index.RootIfd

	getString := func(ifd *exif.Ifd, name string) string {
		results, err := ifd.FindTagWithName(name)
		if err != nil || len(results) == 0 {
			return ""
		}
		val, err := results[0].Value()
		if err != nil {
			return ""
		}
		if s, ok := val.(string); ok {
			return s
		}
		return ""
	}

	out.Make = getString(rootIfd, "Make")
	out.Model = getString(rootIfd, "Model")

	if results, err := rootIfd.FindTagWithName("Orientation"); err == nil && len(results) > 0 {
		if val, err := results[0].Value(); err == nil {
			if b, ok := val.([]uint16); ok && len(b) > 0 {
				o := int(b[0])
				out.Orientation = &o
			}
		}
	}

	if exifIfd, err := rootIfd.ChildWithIfdPath(exifcommon.IfdExifStandardIfdIdentity); err == nil {
		if s := getString(exifIfd, "DateTimeOriginal"); s != "" {
			if t, err := time.Parse("2006:01:02 15:04:05", s); err == nil {
				out.DateTimeOriginal = &t
			}
		}
	}

	if gpsIfd, err := rootIfd.ChildWithIfdPath(exifcommon.IfdGpsInfoStandardIfdIdentity); err == nil {
		if gpsInfo, err := gpsIfd.GpsInfo(); err == nil {
			out.GPS = &domain.GPSData{
				Latitude:  gpsInfo.Latitude.Decimal(),
				Longitude: gpsInfo.Longitude.Decimal(),
			}
		}
	}

	return out, nil
}
