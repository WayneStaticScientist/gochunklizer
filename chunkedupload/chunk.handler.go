package chunkedupload

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (chuck *ChunkUploader) Upload(c *fiber.Ctx) error {
	fileChunk, err := c.FormFile("chunk")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	// chunkIndexStr := c.FormValue("index")
	// chunkIndex, _ := strconv.Atoi(chunkIndexStr)
	fileName := c.FormValue("fileName")
	tempFilePath := fmt.Sprintf("./uploads/temp_%s", strings.ReplaceAll(fileName, " ", "_"))
	f, err := os.OpenFile(tempFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer f.Close()
	chunkData, err := fileChunk.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer chunkData.Close()
	_, err = f.ReadFrom(chunkData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(fiber.Map{"message": "File uploaded successfully"})

}
