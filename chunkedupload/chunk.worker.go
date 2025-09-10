package chunkedupload

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

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
	objectKey := "uploaded_file.zip"

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

	// Create a new S3 client and directly provide the R2 endpoint.
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

	log.Printf("Uploading file: %s to bucket: %s\n", filePath, bucketName)
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		Body:        fileReader.reader,
		ContentType: aws.String("application/octet-stream"), // Set the correct MIME type
	})
	if err != nil {
		log.Fatalf("failed to upload object, %v", err)
	}
	log.Println("File uploaded successfully!")
}
