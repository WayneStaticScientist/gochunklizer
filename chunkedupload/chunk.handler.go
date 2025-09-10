package chunkedupload

import (
	"fmt"
	"log"
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
		uploadDir := "./uploads"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			fmt.Printf("Failed to create upload directory: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to upload file"})
		}
		tempFilePath := fmt.Sprintf("./uploads/temp_%s", strings.ReplaceAll(fileName, " ", "_"))
		chunkCache[fileName] = types.ChunkCache{
			CurrentIndex: int64(chunkIndex),
			TotalChunks:  int64(totalChunks),
			ChunkPath:    tempFilePath,
			Step:         0,
		}
		//user verification here
	}
	if chunkCache[fileName].CurrentIndex != int64(chunkIndex) {
		return c.Status(fiber.StatusBadRequest).SendString("Chunk index mismatch")
	}

	f, err := os.OpenFile(chunkCache[fileName].ChunkPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	defer f.Close()
	chunkData, err := fileChunk.Open()
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	defer chunkData.Close()
	_, err = f.ReadFrom(chunkData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	updatedChunkCache := chunkCachedData
	if chunkIndex == totalChunks-1 {
		updatedChunkCache.Step = 1
	}
	updatedChunkCache.CurrentIndex = chunkCachedData.CurrentIndex + 1
	chunkCache[fileName] = updatedChunkCache
	chunkIndex++
	return c.JSON(fiber.Map{"message": "File uploaded successfully", "progress": chunkIndex / totalChunks})
}
