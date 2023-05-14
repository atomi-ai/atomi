package main

import (
	"github.com/atomi-ai/atomi/utils"
	"log"
)

func main() {
	azureContainerURL := "https://atomidrone.blob.core.windows.net/testing?sp=racwdl&st=2023-05-01T16:59:44Z&se=2024-10-04T00:59:44Z&spr=https&sv=2022-11-02&sr=c&sig=NYvndFrlEpk1cvphkWLI9UEVlLIO2bghCqYVEWXwHyA%3D"
	localFilePath := "/tmp/azure.exp"

	abs, err := utils.NewAzureBlobStorage(azureContainerURL)
	if err != nil {
		log.Fatalf("Failed to connect blob storage: #{err}")
	}

	blobURL, err := abs.UploadFile(localFilePath)
	if err != nil {
		log.Fatalf("Failed to upload file to Azure Blob Storage: %v", err)
	}

	log.Printf("File uploaded successfully. Blob URL: %s\n", blobURL)
}
