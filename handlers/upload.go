package handlers

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	
	_ "image/jpeg"
	_ "image/png"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

// UploadResponse represents the response from image upload
type UploadResponse struct {
	URL       string `json:"url"`
	PublicID  string `json:"public_id"`
	SecureURL string `json:"secure_url"`
}

// UploadImage handles image uploads to Cloudinary
func UploadImage(c *gin.Context) {
	// Parse multipart form (10 MB max)
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to parse form"})
		return
	}

	// Get folder from query parameter or form
	folder := c.Query("folder")
	if folder == "" {
		folder = c.PostForm("folder")
	}
	if folder == "" {
		folder = "playtz"
	}

	// Get file from form
	file, handler, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "No image file provided"})
		return
	}
	defer file.Close()

	// Validate file type
	ext := filepath.Ext(handler.Filename)
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowedExts[ext] {
		c.JSON(400, gin.H{"error": "Invalid file type. Only images are allowed"})
		return
	}

	// Read file into memory
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read file"})
		return
	}

	// Convert to WebP
	fileBytes, err = convertToWebP(fileBytes, ext)
	if err != nil {
		c.JSON(500, gin.H{"error": "Image processing failed: " + err.Error()})
		return
	}

	// Initialize Cloudinary
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")
	
	if cloudName == "" || apiKey == "" || apiSecret == "" {
		c.JSON(500, gin.H{
			"error": "Cloudinary configuration error: missing credentials",
			"debug": map[string]interface{}{
				"cloud_name_set": cloudName != "",
				"api_key_set":    apiKey != "",
				"api_secret_set": apiSecret != "",
			},
		})
		return
	}
	
	// Trim whitespace from credentials
	cloudName = strings.TrimSpace(cloudName)
	apiKey = strings.TrimSpace(apiKey)
	apiSecret = strings.TrimSpace(apiSecret)
	
	// Debug: log what we're sending (first few chars only for security)
	secretPreview := ""
	if len(apiSecret) > 10 {
		secretPreview = apiSecret[:10] + "..."
	} else {
		secretPreview = apiSecret
	}
	
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Upload failed: " + err.Error(),
			"debug": map[string]interface{}{
				"cloud_name": cloudName,
				"api_key": apiKey,
				"api_secret_length": len(apiSecret),
				"api_secret_preview": secretPreview,
			},
		})
		return
	}

	// Generate unique filename with .webp extension
	baseName := filepath.Base(handler.Filename[:len(handler.Filename)-len(ext)])
	filename := fmt.Sprintf("%s_%d.webp", baseName, time.Now().Unix())
	uploadPath := fmt.Sprintf("%s/%s", folder, filename)

	// Create a new reader from the processed bytes
	fileReader := bytes.NewReader(fileBytes)

	// Upload to Cloudinary with WebP format conversion
	ctx := c.Request.Context()
	
	// Try upload without specifying both PublicID and Folder (Cloudinary may handle this better)
	uploadResult, err := cld.Upload.Upload(ctx, fileReader, uploader.UploadParams{
		PublicID:       uploadPath,
		Format:         "webp",
		Transformation: "q_auto:good",
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Upload failed: " + err.Error()})
		return
	}

	// Check if upload result is valid
	if uploadResult == nil {
		c.JSON(500, gin.H{"error": "Upload failed: empty response from Cloudinary"})
		return
	}

	// Check for empty URLs - this indicates a Cloudinary API issue
	if uploadResult.SecureURL == "" && uploadResult.URL == "" {
		c.JSON(500, gin.H{
			"error": "Upload failed: Cloudinary returned empty URLs. This usually indicates invalid credentials or account restrictions.",
			"debug": map[string]interface{}{
				"public_id": uploadResult.PublicID,
				"has_url":   uploadResult.URL != "",
				"has_secure_url": uploadResult.SecureURL != "",
				"result_type": fmt.Sprintf("%T", uploadResult),
			},
		})
		return
	}

	// Return success response
	response := UploadResponse{
		URL:       uploadResult.URL,
		PublicID:  uploadResult.PublicID,
		SecureURL: uploadResult.SecureURL,
	}

	c.JSON(200, response)
}

// UploadMultipleImages handles multiple image uploads
func UploadMultipleImages(c *gin.Context) {
	err := c.Request.ParseMultipartForm(50 << 20) // 50 MB max
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to parse form"})
		return
	}

	folder := c.Query("folder")
	if folder == "" {
		folder = c.PostForm("folder")
	}
	if folder == "" {
		folder = "playtz"
	}

	// Get all files
	form := c.Request.MultipartForm
	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(400, gin.H{"error": "No images provided"})
		return
	}

	// Initialize Cloudinary
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")
	
	if cloudName == "" || apiKey == "" || apiSecret == "" {
		c.JSON(500, gin.H{
			"error": "Cloudinary configuration error: missing credentials",
			"debug": map[string]interface{}{
				"cloud_name_set": cloudName != "",
				"api_key_set":    apiKey != "",
				"api_secret_set": apiSecret != "",
			},
		})
		return
	}
	
	// Trim whitespace from credentials
	cloudName = strings.TrimSpace(cloudName)
	apiKey = strings.TrimSpace(apiKey)
	apiSecret = strings.TrimSpace(apiSecret)
	
	// Debug: log what we're sending (first few chars only for security)
	secretPreview := ""
	if len(apiSecret) > 10 {
		secretPreview = apiSecret[:10] + "..."
	} else {
		secretPreview = apiSecret
	}
	
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Upload failed: " + err.Error(),
			"debug": map[string]interface{}{
				"cloud_name": cloudName,
				"api_key": apiKey,
				"api_secret_length": len(apiSecret),
				"api_secret_preview": secretPreview,
			},
		})
		return
	}

	var responses []UploadResponse
	ctx := c.Request.Context()

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}

		fileBytes, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			continue
		}

		ext := filepath.Ext(fileHeader.Filename)

		// Convert to WebP
		fileBytes, err = convertToWebP(fileBytes, ext)
		if err != nil {
			continue
		}

		baseName := filepath.Base(fileHeader.Filename[:len(fileHeader.Filename)-len(ext)])
		filename := fmt.Sprintf("%s_%d.webp", baseName, time.Now().Unix())
		uploadPath := fmt.Sprintf("%s/%s", folder, filename)

		// Create a new reader from the processed bytes
		fileReader := bytes.NewReader(fileBytes)

		uploadResult, err := cld.Upload.Upload(ctx, fileReader, uploader.UploadParams{
			PublicID:       uploadPath,
			Folder:         folder,
			Format:         "webp",
			Transformation: "q_auto:good",
		})

		if err == nil {
			responses = append(responses, UploadResponse{
				URL:       uploadResult.URL,
				PublicID:  uploadResult.PublicID,
				SecureURL: uploadResult.SecureURL,
			})
		}
	}

	if len(responses) == 0 {
		c.JSON(500, gin.H{"error": "Failed to upload images"})
		return
	}

	c.JSON(200, gin.H{
		"images": responses,
		"count":  len(responses),
	})
}

// convertToWebP optimizes and resizes images before upload
func convertToWebP(imageBytes []byte, originalExt string) ([]byte, error) {
	// Decode image
	img, format, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	// If already WebP and small, return as is
	if format == "webp" && len(imageBytes) < 500000 {
		return imageBytes, nil
	}

	// Resize if too large (max 2000px on longest side for optimization)
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	maxDimension := 2000

	if width > maxDimension || height > maxDimension {
		if width > height {
			img = imaging.Resize(img, maxDimension, 0, imaging.Lanczos)
		} else {
			img = imaging.Resize(img, 0, maxDimension, imaging.Lanczos)
		}
	}

	// Re-encode as JPEG with good quality (Cloudinary will convert to WebP)
	var buf bytes.Buffer
	err = imaging.Encode(&buf, img, imaging.JPEG, imaging.JPEGQuality(85))
	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %v", err)
	}

	return buf.Bytes(), nil
}
