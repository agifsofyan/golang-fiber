package controllers

import (
	"bufio"
	"example/gorest/utils"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type formFile struct {
	File string `json:"file,omitempty" bson:"file,omitempty"`
}

// @Summary     Base64 ENcode file
// @ID 			file.upload
// @Tags        files
// @Accept  	multipart/form-data
// @Produce     json
// @Param   	file formData file true  "Encode the uploaded file"
// @Router      /files [post]
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
	input := new(bson.M)

	if err := c.BodyParser(input); err != nil {
		return utils.FailResponse(c, fiber.StatusBadRequest, "Failed to parse body")
	}

	form := fmt.Sprint(input) // convert to string

	err, path, fileName := utils.WriteBase64ToFile(form)

	if err != nil {
		return utils.FailResponse(c, 400, string(err.Error()))
	}

	relativePath := fmt.Sprintf("%s/%s", path, fileName)

	thumb, err := utils.CreateThumbnail(path, fileName)

	if err != nil {
		return utils.FailResponse(c, 400, string(err.Error()))
	}

	return utils.SuccessResponse(c, fiber.Map{"upload": relativePath, "thumb": thumb}, false)
}
