package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Config struct {
	AlbumURL string `json:"album_url"`
	Port     int    `json:"port"`
}

type ImageInfo struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
}

var (
	config     *Config
	imageCache []ImageInfo
)

func main() {
	// Load configuration
	loadConfig()

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Serve static files
	r.Static("/static", "./")

	// Serve index.html on root
	r.GET("/", func(c *gin.Context) {
		c.File("index.html")
	})

	// API routes under /api
	api := r.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Google Photos Album Random Image API",
				"endpoints": []string{
					"GET /api/random - Get a random image URL from the album",
					"GET /api/img.png - Get a random image file (serves actual image)",
					"GET /api/images - Get all images from the album",
					"GET /api/refresh - Refresh the image cache",
				},
				"features": []string{
					"Auto-refresh cache every 12 hours",
					"CORS enabled for web access",
				},
			})
		})

		api.GET("/random", handleRandomImage)
		api.GET("/img.png", handleRandomImageFile)
		api.GET("/images", handleListImages)
		api.GET("/refresh", handleRefreshCache)
	}

	// Initialize image cache
	refreshImageCache()

	// Start automatic cache refresh every 12 hours
	go startAutoRefresh()

	log.Printf("Server starting on port %d", config.Port)
	r.Run(fmt.Sprintf(":%d", config.Port))
}

func loadConfig() {
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal("Error reading config.json:", err)
	}

	config = &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		log.Fatal("Error parsing config.json:", err)
	}

	if config.AlbumURL == "" {
		log.Fatal("Album URL not found in config.json")
	}

	if config.Port == 0 {
		config.Port = 8080
	}
}

func handleRandomImage(c *gin.Context) {
	if len(imageCache) == 0 {
		c.JSON(404, gin.H{"error": "No images found in album"})
		return
	}

	// Select random image
	randomIndex := rand.Intn(len(imageCache))
	randomImage := imageCache[randomIndex]

	c.JSON(200, randomImage)
}

func handleListImages(c *gin.Context) {
	c.JSON(200, gin.H{
		"images": imageCache,
		"count":  len(imageCache),
	})
}

func handleRefreshCache(c *gin.Context) {
	refreshImageCache()
	c.JSON(200, gin.H{
		"message": "Image cache refreshed",
		"count":   len(imageCache),
	})
}

func handleRandomImageFile(c *gin.Context) {
	if len(imageCache) == 0 {
		c.JSON(404, gin.H{"error": "No images found in album"})
		return
	}

	// Select random image
	randomIndex := rand.Intn(len(imageCache))
	randomImage := imageCache[randomIndex]

	// Fetch the actual image
	resp, err := http.Get(randomImage.URL)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch image: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	// Set appropriate headers
	c.Header("Content-Type", resp.Header.Get("Content-Type"))
	c.Header("Content-Length", resp.Header.Get("Content-Length"))
	c.Header("Cache-Control", "public, max-age=3600") // Cache for 1 hour

	// Stream the image data
	io.Copy(c.Writer, resp.Body)
}

func refreshImageCache() {
	log.Println("Refreshing image cache...")
	
	// Extract album ID from URL
	albumID := extractAlbumID(config.AlbumURL)
	if albumID == "" {
		log.Println("Could not extract album ID from URL")
		return
	}

	// Fetch album page
	albumURL := fmt.Sprintf("https://photos.app.goo.gl/%s", albumID)
	resp, err := http.Get(albumURL)
	if err != nil {
		log.Printf("Error fetching album: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return
	}

	// Extract image URLs from the page
	images := extractImageURLs(string(body))
	
	// Update cache
	imageCache = images
	log.Printf("Cached %d images", len(imageCache))
}

func extractAlbumID(url string) string {
	// Extract album ID from various Google Photos URL formats
	patterns := []string{
		`photos\.app\.goo\.gl/([a-zA-Z0-9_-]+)`,
		`goo\.gl/photos/([a-zA-Z0-9_-]+)`,
		`google\.com/photos/album/([a-zA-Z0-9_-]+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	return ""
}

func extractImageURLs(htmlContent string) []ImageInfo {
	var images []ImageInfo

	// Look for Google Photos image URLs in the HTML
	// These patterns match various Google Photos URL formats
	patterns := []string{
		`https://lh3\.googleusercontent\.com/[^"'\s]+`,
		`https://photos\.google\.com/share/[^"'\s]+`,
		`https://drive\.google\.com/file/d/[^"'\s]+`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllString(htmlContent, -1)
		
		for _, match := range matches {
			// Clean up the URL
			url := strings.TrimSpace(match)
			url = strings.Trim(url, `"'`)
			
			// Extract filename from URL
			filename := extractFilename(url)
			
			images = append(images, ImageInfo{
				URL:      url,
				Filename: filename,
			})
		}
	}

	// Remove duplicates
	uniqueImages := make([]ImageInfo, 0)
	seen := make(map[string]bool)
	
	for _, img := range images {
		if !seen[img.URL] {
			seen[img.URL] = true
			uniqueImages = append(uniqueImages, img)
		}
	}

	return uniqueImages
}

func extractFilename(url string) string {
	// Extract filename from URL
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		filename := parts[len(parts)-1]
		// Remove query parameters
		if idx := strings.Index(filename, "?"); idx != -1 {
			filename = filename[:idx]
		}
		return filename
	}
	return "unknown.jpg"
}

func init() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
}

func startAutoRefresh() {
	// Create a ticker that ticks every 12 hours
	ticker := time.NewTicker(12 * time.Hour)
	defer ticker.Stop()

	log.Println("Auto-refresh started: cache will refresh every 12 hours")

	for {
		select {
		case <-ticker.C:
			log.Println("Auto-refreshing image cache...")
			refreshImageCache()
		}
	}
} 