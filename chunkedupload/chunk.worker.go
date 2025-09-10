package chunkedupload

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"opechains.shop/chunklizer/v2/types"
)

type fileWithProgress struct {
	reader io.Reader
	size   int64
}

func (chunk *ChunkUploader) CleanUpTemp() {
	for {
		time.Sleep(time.Minute * 30)
		chunkCacheMutex.Lock()
		for k, v := range chunkCache {
			if (v.LastAccess - time.Now().Unix()) > 300 {
				os.Remove(v.ChunkPath)
				delete(chunkCache, k)
			}
		}
		chunkCacheMutex.Unlock()
	}
}

func (chunk *ChunkUploader) Work() {
	for v := range chunkChan {
		chunk.UploadToCloud(v)
	}
}

func (c *ChunkUploader) UploadToCloud(chunk types.ChunkCache) {
	accountID := os.Getenv("CACCOUNT_ID")
	bucketName := "pictures"
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	filePath := chunk.ChunkPath
	if accountID == "" || bucketName == "" || accessKeyID == "" || secretAccessKey == "" {
		log.Fatal("Error: Missing R2 environment variables.")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("auto"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
	)
	if err != nil {
		log.Fatalf("Failed to load SDK configuration: %v", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID))
	})
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file, %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("failed to get file info, %v", err)
	}

	fileReader := &fileWithProgress{
		reader: file,
		size:   fileInfo.Size(),
	}

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Body:        fileReader.reader,
		Bucket:      aws.String(bucketName),
		Key:         aws.String(chunk.FileName),
		ContentType: aws.String("application/octet-stream"),
	})
	if err != nil {
		log.Fatalf("failed to upload object, %v", err)
	}
	log.Println("File uploaded successfully!")
	err = os.Remove(filePath)
	if err != nil {
		log.Printf("Failed to delete local file %s: %v", filePath, err)
	} else {
		log.Printf("Successfully deleted local file %s", filePath)
	}

}
