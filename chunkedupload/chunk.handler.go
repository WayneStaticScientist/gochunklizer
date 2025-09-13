package chunkedupload

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"opechains.shop/chunklizer/v2/types"
	"opechains.shop/chunklizer/v2/user"
)

// ->[post] /utils/files?t=
func (chuck *ChunkUploader) Upload(c *fiber.Ctx) error {
	objectId := c.Get("X-Object-Id")
	userToken := c.Get("X-Upload-Id")
	if strings.Trim(userToken, " ") == "" || strings.Trim(objectId, " ") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}
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
		log.Println(errtotal.Error())
		return c.Status(fiber.StatusBadRequest).SendString("Invalid chunk index")
	}
	fileName := c.FormValue("fileName")
	fileType := c.FormValue("fileType")
	if strings.Trim(fileType, " ") == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid File Type")
	}
	if _, ok := chunkCache[userToken]; !ok {
		if err := user.VerifyToken(userToken); err != nil {
			log.Println("Error from third server ", err.Error())
			return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
		}
		baseOnCurrentTime := strconv.FormatInt(c.Context().Time().Unix(), 10)
		uploadDir := fmt.Sprintf("./uploads/temp/%s/%s", baseOnCurrentTime, objectId)
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			fmt.Printf("Failed to create upload directory: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to upload file"})
		}
		fileName = strings.ReplaceAll(fileName, "..", "")
		fileName = strings.ReplaceAll(fileName, "/", "")
		fileName = strings.ReplaceAll(fileName, "\\", "")
		fileName = strings.ReplaceAll(fileName, " ", "_")
		if fileName == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid file name")
		}
		tempFilePath := fmt.Sprintf("%s/%s", uploadDir, fileName)
		chunkCache[userToken] = types.ChunkCache{
			Step:         0,
			FileName:     fileName,
			FileType:     fileType,
			ObjectId:     objectId,
			Token:        userToken,
			ChunkPath:    tempFilePath,
			LastAccess:   time.Now().Unix(),
			CurrentIndex: int64(chunkIndex),
			TotalChunks:  int64(totalChunks),
		}
	}
	if chunkCache[userToken].CurrentIndex != int64(chunkIndex) {
		return c.Status(fiber.StatusBadRequest).SendString("Chunk index mismatch")
	}
	f, err := os.OpenFile(chunkCache[userToken].ChunkPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	defer f.Close()
	chunkData, err := fileChunk.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	defer chunkData.Close()
	_, err = f.ReadFrom(chunkData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	updatedChunkCache := chunkCache[userToken]
	if chunkIndex == totalChunks-1 {
		chunkChan <- updatedChunkCache
		delete(chunkCache, userToken)
		return c.JSON(fiber.Map{"message": "File uploaded successfully", "progress": 1})
	}
	updatedChunkCache.ChunkPath = chunkCache[userToken].ChunkPath
	updatedChunkCache.TotalChunks = chunkCache[userToken].TotalChunks
	updatedChunkCache.CurrentIndex = chunkCache[userToken].CurrentIndex + 1
	updatedChunkCache.LastAccess = time.Now().Unix()
	chunkCache[userToken] = updatedChunkCache
	return c.JSON(fiber.Map{"message": "File uploaded successfully", "progress": chunkIndex / totalChunks})
}

func (chuck *ChunkUploader) Update(c *fiber.Ctx) error {
	objectId := c.Get("X-Object-Id")
	userToken := c.Get("X-Upload-Id")
	objectKey := c.Get("X-Object-Key", "")
	if strings.Trim(userToken, " ") == "" || strings.Trim(objectId, " ") == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}
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
		log.Println(errtotal.Error())
		return c.Status(fiber.StatusBadRequest).SendString("Invalid chunk index")
	}
	fileName := c.FormValue("fileName")
	fileType := c.FormValue("fileType")
	if strings.Trim(fileType, " ") == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid File Type")
	}
	if _, ok := chunkCache[userToken]; !ok {
		if err := user.VerifyToken(userToken); err != nil {
			log.Println("Error from third server ", err.Error())
			return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
		}
		if strings.ReplaceAll(objectKey, " ", "") != "" {
			if err := chuck.DeleteFromCloud(c.Context(), objectKey); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to replace file"})
			}
		}
		baseOnCurrentTime := strconv.FormatInt(c.Context().Time().Unix(), 10)
		uploadDir := fmt.Sprintf("./uploads/temp/%s/%s", baseOnCurrentTime, objectId)
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			fmt.Printf("Failed to create upload directory: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to upload file"})
		}
		fileName = strings.ReplaceAll(fileName, "..", "")
		fileName = strings.ReplaceAll(fileName, "/", "")
		fileName = strings.ReplaceAll(fileName, "\\", "")
		fileName = strings.ReplaceAll(fileName, " ", "_")
		if fileName == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid file name")
		}
		tempFilePath := fmt.Sprintf("%s/%s", uploadDir, fileName)
		chunkCache[userToken] = types.ChunkCache{
			Step:         0,
			FileName:     fileName,
			FileType:     fileType,
			ObjectId:     objectId,
			Token:        userToken,
			ChunkPath:    tempFilePath,
			LastAccess:   time.Now().Unix(),
			CurrentIndex: int64(chunkIndex),
			TotalChunks:  int64(totalChunks),
		}
	}
	if chunkCache[userToken].CurrentIndex != int64(chunkIndex) {
		return c.Status(fiber.StatusBadRequest).SendString("Chunk index mismatch")
	}
	f, err := os.OpenFile(chunkCache[userToken].ChunkPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	defer f.Close()
	chunkData, err := fileChunk.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	defer chunkData.Close()
	_, err = f.ReadFrom(chunkData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	updatedChunkCache := chunkCache[userToken]
	if chunkIndex == totalChunks-1 {
		chunkChan <- updatedChunkCache
		delete(chunkCache, userToken)
		return c.JSON(fiber.Map{"message": "File uploaded successfully", "progress": 1})
	}
	updatedChunkCache.ChunkPath = chunkCache[userToken].ChunkPath
	updatedChunkCache.TotalChunks = chunkCache[userToken].TotalChunks
	updatedChunkCache.CurrentIndex = chunkCache[userToken].CurrentIndex + 1
	updatedChunkCache.LastAccess = time.Now().Unix()
	chunkCache[userToken] = updatedChunkCache
	return c.JSON(fiber.Map{"message": "File uploaded successfully", "progress": chunkIndex / totalChunks})
}

func (c *ChunkUploader) RequestDeleteFile(ctx *fiber.Ctx) error {
	objectKey := ctx.Get("X-Object-Key", "")
	if os.Getenv("S2S_API_KEY") != ctx.Get("Authorization") {
		return ctx.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}
	accountID := os.Getenv("R2_ACCOUNT_ID")
	bucketName := os.Getenv("R2_BUCKET_NAME")
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")

	if accountID == "" || bucketName == "" || accessKeyID == "" || secretAccessKey == "" {
		log.Fatal("Error: Missing R2 environment variables. Please check your .env file or system settings.")
	}
	if strings.Trim(objectKey, " ") == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid object key")
	}

	if err := c.DeleteFromCloud(ctx.Context(), objectKey); err != nil {
		log.Printf("Failed to delete object: %v", err)
		ctx.Status(fiber.StatusInternalServerError).SendString("Failed to delete object")
	}
	log.Println("Object deleted successfully! 🗑️")
	return ctx.JSON(fiber.Map{"message": "File deleted successfully"})
}
