# Google Photos Album Random Image API

A simple Go API that serves random images from a public Google Photos album.

## Setup

1. **Install Go** (if not already installed)
2. **Update config.json** with your public Google Photos album URL:

```json
{
  "album_url": "https://photos.app.goo.gl/YOUR_ALBUM_ID",
  "port": 8080
}
```

3. **Run the API**:
```bash
go mod tidy
go run main.go
```

## API Endpoints

- `GET /` - Demo page (serves index.html)
- `GET /api` - API information
- `GET /api/random` - Get a random image URL from the album
- `GET /api/img.png` - Get a random image file (serves actual image)
- `GET /api/images` - Get all images from the album
- `GET /api/refresh` - Refresh the image cache

## Example Usage

```bash
# Visit the demo page
open http://localhost:8080

# Get a random image URL
curl http://localhost:8080/api/random

# Get a random image file (downloads actual image)
curl http://localhost:8080/api/img.png -o random_image.png

# Get all images
curl http://localhost:8080/api/images

# Refresh cache
curl http://localhost:8080/api/refresh
```

## How it works

1. The API extracts the album ID from the provided Google Photos URL
2. It fetches the album page and extracts image URLs
3. Images are cached in memory for fast access
4. The `/random` endpoint returns a random image from the cache
5. **Auto-refresh**: Cache automatically refreshes every 12 hours to keep images up-to-date

## Requirements

- Public Google Photos album URL
- Go 1.21+
- Internet connection to fetch album data 