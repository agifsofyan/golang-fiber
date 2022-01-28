package controllers

import (
	"bufio"
	"encoding/base64"
	"example/gorest/utils"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func FileEncode(c *fiber.Ctx) error {
	file, msg := utils.Upload(c, "storage/temp", "file")

	if file == nil {
		return utils.FailResponse(c, 400, msg)
	}

	f, _ := os.Open(file["fullPath"].(string))
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	imgUrl := utils.ToBase64(content)

	return utils.SuccessResponse(c, fiber.Map{"imgUrl": imgUrl}, false)
}

func FileDecode(c *fiber.Ctx) error {
	// file := c.FormValue("file")

	input := new(bson.M)

	if err := c.BodyParser(input); err != nil {
		return utils.FailResponse(c, fiber.StatusBadRequest, "Failed to parse body")
	}

	form := fmt.Sprint(input) // convert to string

	typeFile, strFile := utils.FromBase64(form)

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(strFile))

	m, formatSring, err := image.Decode(reader)

	if err != nil {
		return utils.FailResponse(c, 400, string(err.Error()))
	}

	bounds := m.Bounds()
	log.Println(bounds, "<->", formatSring)

	pngFilename := "./storage/upload/test." + typeFile
	f, err := os.OpenFile(pngFilename, os.O_WRONLY|os.O_CREATE, 07777)

	if err != nil {
		return utils.FailResponse(c, 400, string(err.Error()))
	}

	err = png.Encode(f, m)

	if err != nil {
		return utils.FailResponse(c, 400, string(err.Error()))
	}

	return utils.SuccessResponse(c, fiber.Map{"pngFilename": pngFilename}, false)
}
