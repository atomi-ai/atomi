package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type ImageController interface {
	UploadImage(c *gin.Context)
}

type ImageControllerImpl struct {
	blobStorage utils.BlobStorage
}

func NewImageController(blobStorage utils.BlobStorage) ImageController {
	return &ImageControllerImpl{blobStorage: blobStorage}
}

func (ic *ImageControllerImpl) UploadImage(c *gin.Context) {
	// 限制只有manager可以上传图片。
	user, _ := c.Get("user")
	manager := user.(*models.User)

	if !isAuthorized(manager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide an image file"})
		return
	}

	now := time.Now()
	formattedTime := now.Format("20060102-150405")
	newFileName := fmt.Sprintf("%s-%s", formattedTime, file.Filename)
	tempFilePath := filepath.Join(os.TempDir(), newFileName)

	err = c.SaveUploadedFile(file, tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file"})
		return
	}
	log.Debugf("xfguo: temporarily saved the file in '%v'", tempFilePath)

	uploadedFileURL, err := ic.blobStorage.UploadFile(tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to Azure Blob Storage"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"uploaded_file_url": uploadedFileURL})
}
