package handlers

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	env "fictsu_backend/config"

	gsc "cloud.google.com/go/storage"
)

func UploadImageToFirebase(file multipart.File, fileHeader *multipart.FileHeader, objectPath string, bucketName string) (string, error) {
	ctx := context.Background()

	storageClient, err := env.FirebaseApp.Storage(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get Firebase storage client: %v", err)
	}
	bucket, err := storageClient.Bucket(bucketName)
	if err != nil {
		return "", fmt.Errorf("failed to get default bucket: %v", err)
	}
	writer := bucket.Object(objectPath).NewWriter(ctx)
	writer.ContentType = fileHeader.Header.Get("Content-Type")
	if _, err := io.Copy(writer, file); err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %v", err)
	}
	if err := bucket.Object(objectPath).ACL().Set(ctx, gsc.AllUsers, gsc.RoleReader); err != nil {
		return "", fmt.Errorf("failed to make file public: %v", err)
	}
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectPath)
	return publicURL, nil
}
