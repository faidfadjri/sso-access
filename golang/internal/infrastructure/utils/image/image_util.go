package image

import (
	pkgErrors "akastra-access/internal/pkg/errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nfnt/resize"
)

func ProcessImage(file multipart.File, header *multipart.FileHeader, destDir string) (string, error) {
	// Validate file format
	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	fileType := http.DetectContentType(buff)
	if fileType != "image/jpeg" && fileType != "image/png" {
		return "", pkgErrors.ErrInvalidFormat
	}

	// Reset file pointer
	if _, err := file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("failed to reset file pointer: %v", err)
	}

	// Decode image
	img, format, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %v", err)
	}

	// Resize image (width 300, preserve aspect ratio)
	m := resize.Resize(300, 0, img, resize.Lanczos3)

	// Generate filename
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	
	// Create directory if not exists
	uploadDir := destDir
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	outPath := filepath.Join(uploadDir, filename)
	out, err := os.Create(outPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	// Encode and save (compress)
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(out, m, &jpeg.Options{Quality: 75})
	case "png":
		err = png.Encode(out, m)
	default:
		// Default to jpeg for others
		err = jpeg.Encode(out, m, &jpeg.Options{Quality: 75})
	}

	if err != nil {
		return "", fmt.Errorf("failed to encode image: %v", err)
	}

	return "/" + outPath, nil
}
