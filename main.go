package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/justikun/metadata-viewer/pkg/jpg"
	"github.com/justikun/metadata-viewer/pkg/metadata"
)

func main() {
	//photos := [3]string{"test-photos/test-image-1.jpeg", "test-photos/test-image-2.tiff", "test-photos/test-image-3"}

	images, err := GetImageFiles("test-photos")
	if err != nil {
		fmt.Println(err)
	}

	for _, image := range images {
		images, err := parseImgData(image.ImagePath)
		if err != nil {
			fmt.Printf("Failed to parse image at %s", images.ImagePath)
		}
	}

	print(images[0].ImagePath)
}

func parseImgData(imgPath string) (metadata.ImageData, error) {
	imgData := metadata.ImageData{}
	imgData.ImagePath = imgPath

	file, err := os.Open(imgPath)
	if err != nil {
		return imgData, fmt.Errorf("error: %v, is NOT a valid image path\n", imgPath)
	}
	defer file.Close()

	for {
		// find marker
		marker := make([]byte, 2)
		_, err = file.Read(marker)
		if err != nil {
			return imgData, err
		}

		if marker[0] != 0xFF {
			return imgData, fmt.Errorf("Hex: %x. Not a valid marker\n", marker[0])
		}

		switch marker[1] {
		case 0xD8: // SOI - start of image
			println("case soi:")
		case 0xD9: // EOI - end of image
			return imgData, fmt.Errorf("case eoi")
		case 0xE0: // APP0 - jfif marker
			println("case app0:")
			jpg.ParseAPP0(file, &imgData)
			if err != nil {
				fmt.Printf("Failed to parse APPO at %s\n", imgPath)
			}
		case 0xE1: // APP1
			println("case app1:")
			err := jpg.ParseAPP1(file, &imgData)
			if err != nil {
				fmt.Printf("Failed to parse APP1 at %s\n", imgPath)
			}
			println("APP1 DONE")
			helperPrintData(imgData)
			println("--------------")

			return imgData, nil
		case 0xDA: // SOS - image stream
			return imgData, nil
		default:
			return imgData, fmt.Errorf("No marker found")
		}
	}
}

func GetImageFiles(dirPath string) ([]metadata.ImageData, error) {
	var imageFiles []metadata.ImageData

	allowedExtensions := map[string]struct{}{
		".jpg":  {},
		".jpeg": {},
		".png":  {},
		".tiff": {},
		".tif":  {},
		".webp": {},
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return []metadata.ImageData{}, fmt.Errorf("Failed to read directory %s: %w", dirPath, err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		fileName := entry.Name()
		ext := strings.ToLower(filepath.Ext(fileName))
		_, ok := allowedExtensions[ext]
		if ok {
			// construct absolute path
			absPath := filepath.Join(dirPath, fileName)
			image := metadata.ImageData{ImagePath: absPath, MetaData: metadata.MetaData{}}
			imageFiles = append(imageFiles, image)
		}
	}
	return imageFiles, nil
}

func helperPrintData(data metadata.ImageData) {

	fmt.Println("Image: ", data.ImagePath)

	fmt.Println("")
	fmt.Println("-----END OF Print Data-----")
}
