package controllers

import (
	"errors"
	"mime/multipart"
	"net/http"
	"os"
	"smart-serve/models"
	"smart-serve/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const MaxFileSize = 1024 * 1024 // 1MB

func isImage(fileHeader *multipart.FileHeader) error {
	contentType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return errors.New("file is not an image")
	}
	if fileHeader.Size > MaxFileSize {
		return errors.New("file size exceeds 1MB")
	}
	return nil
}

// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 201 {object} Response{data=models.File}
// @Router /files [post]
// @Security BearerAuth
func UploadFile(c *gin.Context) {
	restaurantId := c.GetString("restaurantId")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
		return
	}

	// Check if the file is an image and size is less than 2MB
	if err := isImage(header); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	newFile, err := utils.Uploader.UploadFile(file, header, uuid.New())
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	parsedRestaurantId, _ := uuid.Parse(restaurantId)

	var newFileModel models.File = models.File{
		ID:           newFile.ID,
		Name:         newFile.Name,
		MineType:     newFile.MineType,
		RestaurantID: parsedRestaurantId,
	}

	fileResponse, err := models.CreateFile(newFileModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    fileResponse,
	})
}

// @Tags Files
// @Accept  json
// @Produce  json
// @Param id path string true "File ID"
// @Success 302
// @Router /files/{id} [get]
func GetFile(c *gin.Context) {
	id := c.Param("id")

	c.Redirect(http.StatusFound, os.Getenv("S3_URL")+id)
}

// @Tags Files
// @Accept multipart/form-data
// @Produce  json
// @Param id path string true "File ID"
// @Param file formData file true "File to upload"
// @Success 200 {object} Response{data=models.File}
// @Router /files/{id} [put]
// @Security BearerAuth
func UpdateFile(c *gin.Context) {
	id := c.Param("id")
	restaurantId := c.GetString("restaurantId")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Check if the file is an image and size is less than 2MB
	if err := isImage(header); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oldFile, err := models.GetFile(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if oldFile.RestaurantID.String() != restaurantId {
		c.JSON(http.StatusForbidden, Response{
			Success: false,
			Message: "Access denied",
		})
		return
	}

	key, _ := uuid.Parse(id)
	newFile, err := utils.Uploader.UploadFile(file, header, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	var newFileModel models.File = models.File{
		ID:       newFile.ID,
		Name:     newFile.Name,
		MineType: newFile.MineType,
	}
	updatedFile, err := models.UpdateFile(id, newFileModel)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Update file successfully",
		Data:    updatedFile,
	})
}

// @Tags Files
// @Accept  json
// @Produce  json
// @Param id path string true "File ID"
// @Success 200 {object} Response{data=nil}
// @Router /files/{id} [delete]
// @Security BearerAuth
func DeleteFile(c *gin.Context) {
	id := c.Param("id")
	restaurantId := c.GetString("restaurantId")

	oldFile, err := models.GetFile(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if oldFile.RestaurantID.String() != restaurantId {
		c.JSON(http.StatusForbidden, Response{
			Success: false,
			Message: "Access denied",
		})
		return
	}

	if err := models.DeleteFile(id); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	utils.Uploader.RemoveFile(oldFile.ID)

	// TODO: Remove file from S3 error
	// uuid, _ := uuid.Parse(id)
	// err := utils.Uploader.RemoveFile(uuid)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, Response{
	// 		Success: false,
	// 		Message: err.Error(),
	// 	})
	// 	return
	// }

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "File deleted successfully",
	})
}
