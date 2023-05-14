package main

import (
	"github.com/atomi-ai/atomi/utils"
	"log"
)

func main() {
	azureContainerURL := "https://atomidrone.blob.core.windows.net/images?sv=2022-11-02&ss=bfqt&srt=sco&sp=rwdlacupiytfx&se=2024-04-16T23:57:30Z&st=2023-05-03T15:57:30Z&spr=https&sig=6rYp3%2Fak5JnN96%2BE5BmGDrzUyko1NRNg8xSp6p8CpXw%3D"
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
