package services

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
)

func GetBase64ImageByUrl(imageUrl string) (string, error) {
	resp, err := http.Get(imageUrl)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch image, status code: %d", resp.StatusCode)
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %v", err)
	}

	imageBase64 := base64.StdEncoding.EncodeToString(imageData)
	return fmt.Sprintf("data:image/png;base64,%s", imageBase64), nil
}

func GetImageByURL(imageUrl string) image.Image {
	resp, err := http.Get(imageUrl)
	if err != nil {
		log.Println("Error fetching cover image:", err)
		return nil
	}
	defer resp.Body.Close()

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading cover image data:", err)
		return nil
	}
	contentType := http.DetectContentType(imgData)

	var img image.Image
	switch contentType {
	case "image/jpeg":
		img, err = jpeg.Decode(bytes.NewReader(imgData))
	case "image/png":
		img, err = png.Decode(bytes.NewReader(imgData))
	default:
		log.Printf("Unsupported image type: %s", contentType)
		return nil
	}

	if err != nil {
		log.Printf("Error decoding image: %v", err)
		return nil
	}

	return img
}
