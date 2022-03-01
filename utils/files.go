package utils

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nfnt/resize"
)

func Upload(c *fiber.Ctx, dir string, inputName string) (fiber.Map, string) {
	file, err := c.FormFile(inputName)

	if err != nil {
		return nil, "fail get the file"
	}

	fileName := file.Filename
	contentType := file.Header["Content-Type"][0]
	fileZise := file.Size
	filePath := fmt.Sprintf("./%s/%s", dir, fileName)

	log.Println("path:", filePath)

	err = c.SaveFile(file, filePath)

	if err != nil {
		return nil, "cannot save the file"
	}

	return fiber.Map{
		"name":     fileName,
		"size":     fileZise,
		"type":     contentType,
		"fullPath": filePath,
	}, ""
}

func WriteBase64ToFile(fileString string) (e error, path, fileName string) {
	typeFile, strFile := FromBase64(fileString)

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(strFile))

	m, _, err := image.Decode(reader)

	if err != nil {
		return err, "", ""
	}

	m.Bounds()

	dateNow := time.Now().Unix()
	path = "/storage/upload"
	fileName = fmt.Sprintf("%s.%s", fmt.Sprint(dateNow), typeFile)
	dir := fmt.Sprintf("%s/%s", path, fileName)
	relativePath := "." + dir

	f, err := os.OpenFile(relativePath, os.O_WRONLY|os.O_CREATE, 07777)

	if err != nil {
		return err, "", ""
	}

	err = png.Encode(f, m)

	if err != nil {
		return err, "", ""
	}

	return nil, path, fileName
}

func CreateThumbnail(path, fileName string) (map[string]string, error) {
	oriPath := fmt.Sprintf(".%s/%s", path, fileName)
	thumbPath := make(map[string]string)
	file, err := os.Open(oriPath)
	if err != nil {
		return thumbPath, err
	}

	var img image.Image
	if !strings.Contains(fileName, "jpeg") {
		img, err = png.Decode(file)
	} else {
		img, err = jpeg.Decode(file)
	}

	// decode jpeg into image.Image
	if err != nil {
		return thumbPath, err
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	small := resize.Resize(100, 0, img, resize.Lanczos3)
	medium := resize.Resize(500, 0, img, resize.Lanczos3)

	smallPath := "/storage/_thumb/small/" + fileName
	smallOut, err := os.Create("." + smallPath)
	if err != nil {
		return thumbPath, err
	}
	defer smallOut.Close()

	mediumPath := "/storage/_thumb/medium/" + fileName
	mediumOut, err := os.Create("." + mediumPath)
	if err != nil {
		return thumbPath, err
	}
	defer mediumOut.Close()

	if !strings.Contains(fileName, "jpeg") {
		// write new image to file
		png.Encode(smallOut, small)
		png.Encode(mediumOut, medium)
	} else {
		// write new image to file
		jpeg.Encode(smallOut, small, nil)
		jpeg.Encode(mediumOut, medium, nil)
	}

	thumbPath["smallPath"] = smallPath
	thumbPath["mediumPath"] = mediumPath

	return thumbPath, nil
}
