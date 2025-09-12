package chunkedupload

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"opechains.shop/chunklizer/v2/types"
)

func (chunk *ChunkUploader) CleanUpTemp() {
	for {
		time.Sleep(time.Minute * 30)
		chunkCacheMutex.Lock()
		for k, v := range chunkCache {
			if (time.Now().Unix() - v.LastAccess) > 300 {
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
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	bucketName := getBucketName(strings.ToLower(chunk.FileName))
	filePath := chunk.ChunkPath
	chunk.FileName = fmt.Sprintf("%s_%s", chunk.ObjectId, chunk.FileName)
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
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Body:        file,
		Bucket:      aws.String("assets"),
		Key:         aws.String(fmt.Sprintf("%s%s", bucketName, chunk.FileName)),
		ContentType: aws.String("application/octet-stream"),
	})
	if err != nil {
		log.Fatalf("failed to upload object, %v", err)
	}
	file.Close()
	err = os.Remove(filePath)
	if err != nil {
		log.Printf("Failed to delete local file %s: %v", filePath, err)
	} else {
		log.Printf("Successfully deleted local file %s", filePath)
	}
	trials := 0
	for {
		err := handShakeServer(chunk.Token, fmt.Sprintf("%s%s", bucketName, chunk.FileName), chunk.ObjectId, chunk.FileName)
		if err == nil {
			break
		}
		trials++
		if trials > 10 {
			log.Println("Project with name ", chunk.FileName, "failed to upload")
			break
		}
	}
	os.Remove(chunk.ChunkPath)
}

func getBucketName(s string) string {
	if strings.HasSuffix(s, ".png") || strings.HasSuffix(s, ".jpg") || strings.HasSuffix(s, ".jpeg") || strings.HasSuffix(s, ".gif") {
		return "images/"
	}
	if strings.HasSuffix(s, ".mp4") || strings.HasSuffix(s, ".mov") || strings.HasSuffix(s, ".avi") || strings.HasSuffix(s, ".webm") || strings.HasSuffix(s, ".mkv") {
		return "videos/"
	}
	return "documents/"
}

func handShakeServer(token string, coverUrl string, objectId string, objectKey string) error {
	userData := map[string]any{
		"media":    coverUrl,
		"objectId": objectId,
		"key":      objectKey,
	}

	jsonData, err := json.Marshal(userData)
	if err != nil {
		log.Printf("Error marshaling JSON: %s\n", err)
		return err
	}
	req, err := http.NewRequest("POST", os.Getenv("NEXT_PUBLIC_API_URL")+"/utils/files?t="+token, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error marshaling JSON: %s\n", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %s", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("received non-201 status code: %d", resp.StatusCode)
		return fmt.Errorf("received non-201 status code: %d", resp.StatusCode)
	}
	log.Println("Uploaded to server successfully: ", coverUrl)
	return nil
}

func (chunk *ChunkUploader) DeleteFromCloud(ctx context.Context, key string) error {
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("auto"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
	)
	if err != nil {
		log.Printf("Failed to load SDK configuration: %v", err)
		return fmt.Errorf("failed to load data")
	}
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID))
	})
	_, err = client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("ASSETS_BUCKET_KEY")),
		Key:    aws.String(key),
	})
	return err
}
