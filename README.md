# SnapCheck

SnapCheck is a Go microservice for validating and analyzing JPEG images. It performs binary-level validation, extracts EXIF metadata, and applies quality heuristics.

## Features

- **Binary Validation**: Checks JPEG SOI/EOI markers, file signature, and size limits (2KB - 20MB).
- **EXIF Extraction**: Extracts DateTimeOriginal, Make, Model, Orientation, and GPS data.
- **Heuristics**:
  - Detects timestamp mismatch (EXIF vs. Upload time).
  - Flags missing orientation.
  - Flags low resolution images (< 200px).
- **Graceful Error Handling**: returns structured JSON responses.

## Getting Started

### Prerequisites

- Go 1.20 or later

### Installation

Clone the repository and install dependencies:

```bash
git clone <repository-url>
cd snapcheck
go mod download
```

### Running the Service

Start the server:

```bash
go run cmd/server/main.go
```

The server will start on port `8080`.

## Usage

### Analyze an Image

Use `curl` to upload a JPEG file. You can optionally provide `last_modified` (RFC3339 format) to check for timestamp mismatches.

The repository includes sample images (e.g., `Canon_40D.jpg`) that you can use for testing.

```bash
curl -X POST http://localhost:8080/v1/analyze \
  -F "file=@Canon_40D.jpg" \
  -F "last_modified=2023-10-27T10:00:00Z"
```

### Response Example

```json
{
  "image_info": {
    "width": 1024,
    "height": 768,
    "size_bytes": 450000
  },
  "exif": {
    "has_exif": true,
    "date_time_original": "2023-10-27T10:00:00Z",
    "make": "Canon",
    "model": "Canon EOS R5",
    "orientation": 1,
    "gps": {
      "latitude": 34.0522,
      "longitude": -118.2437
    }
  },
  "heuristics": {
    "timestamp_mismatch": false,
    "orientation_missing": false,
    "low_resolution": false
  }
}
```

### Error Response Example

```json
{
  "image_info": {
    "width": 0,
    "height": 0,
    "size_bytes": 0
  },
  "heuristics": {
    "timestamp_mismatch": false,
    "orientation_missing": false,
    "low_resolution": false
  },
  "error": "file size must be between 2KB and 20MB"
}
```

## Testing

Run unit tests:

```bash
go test ./...
```
