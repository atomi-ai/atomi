package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

type BlobStorage interface {
	UploadFile(filePath string) (string, error)
}

type AzureBlobStorage struct {
	ContainerURL azblob.ContainerURL
}

func NewAzureBlobStorage(containerURL string) (*AzureBlobStorage, error) {
	containerURLParsed, err := url.Parse(containerURL)
	if err != nil {
		return nil, err
	}

	container := azblob.NewContainerURL(*containerURLParsed, azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{}))
	blobStorage := &AzureBlobStorage{ContainerURL: container}
	if err := blobStorage.validateAzureBlobStorage(); err != nil {
		return nil, err
	}
	return blobStorage, nil
}

func (abs *AzureBlobStorage) validateAzureBlobStorage() error {
	// Create a temporary file with content "host_port_timestamp"
	tempFilePath := filepath.Join(os.TempDir(), "init_check.temp")
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return fmt.Errorf("azure blob storage error: failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	timestamp := time.Now().Format(time.RFC3339)
	if _, err := tempFile.WriteString(timestamp); err != nil {
		return fmt.Errorf("azure blob storage error: failed to write to temp file: %w", err)
	}

	uploadedURL, err := abs.UploadFileWithTimeout(tempFile.Name(), 3*time.Second)
	if err != nil {
		return fmt.Errorf("azure blob storage error: failed to upload temp file to Azure Blob Storage: %w", err)
	}

	// Download the temporary file from the container
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	response, err := client.Get(uploadedURL)
	if err != nil {
		return fmt.Errorf("azure blob storage error: failed to download temp file from Azure Blob Storage: %w", err)
	}
	defer response.Body.Close()

	uploadedContentBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("azure blob storage error: failed to read downloaded temp file content: %w", err)
	}

	uploadedContent := string(uploadedContentBytes)

	// Check if the content matches
	if uploadedContent != timestamp {
		return fmt.Errorf("azure blob storage error: content mismatch: expected %s, got %s", timestamp, uploadedContent)
	}
	return nil
}

func (abs *AzureBlobStorage) UploadFile(filePath string) (string, error) {
	return abs.UploadFileWithTimeout(filePath, 10*time.Second)
}

func (abs *AzureBlobStorage) UploadFileWithTimeout(filePath string, timeout time.Duration) (string, error) {
	// Set the default timeout value to 3 seconds
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fileName := filepath.Base(filePath)
	blockBlobURL := abs.ContainerURL.NewBlockBlobURL(fileName)

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = azblob.UploadFileToBlockBlob(ctx, file, blockBlobURL, azblob.UploadToBlockBlobOptions{})
	if err != nil {
		return "", err
	}

	// Get the full URL including SAS token
	fullURL := blockBlobURL.URL()

	// Create a new URL without query parameters (i.e., without SAS token)
	cleanURL := url.URL{
		Scheme: fullURL.Scheme,
		Host:   fullURL.Host,
		Path:   fullURL.Path,
	}

	return cleanURL.String(), nil
}
