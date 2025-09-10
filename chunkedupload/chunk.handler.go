package chunkedupload

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"opechains.shop/chunklizer/v2/types"
)

func (chuck *ChunkUploader) Upload(c *fiber.Ctx) error {
	fileChunk, err := c.FormFile("chunk")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	chunkIndex, errIndex := strconv.Atoi(c.FormValue("index"))
	if errIndex != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid chunk index")
	}
	totalChunks, errtotal := strconv.Atoi(c.FormValue("totalChunks"))
	if errtotal != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid chunk index")
	}
	fileName := c.FormValue("fileName")
	chunkCachedData, ok := chunkCache[fileName]
	if !ok {
		chunkCache[fileName] = types.ChunkCache{
			CurrentIndex: int64(chunkIndex),
			TotalChunks:  int64(totalChunks),
		}
		//user verification here
	}
	if chunkCachedData.CurrentIndex != int64(chunkIndex) {
		return c.Status(fiber.StatusBadRequest).SendString("Chunk index mismatch")
	}
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
	return c.JSON(fiber.Map{"message": "File uploaded successfully", "progress": chunkIndex / totalChunks})

}
