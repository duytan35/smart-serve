package controllers

import (
	"net/http"
	"os"
	"smart-serve/models"
	"smart-serve/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 201 {object} Response{data=models.File}
// @Router /files [post]
// @Security BearerAuth
func UploadFile(c *gin.Context) {
	restaurantId := c.GetString("id")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
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
	newFile.RestaurantID = parsedRestaurantId

	fileResponse, err := models.CreateFile(newFile)
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
	restaurantId := c.GetString("id")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
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

	updatedFile, err := models.UpdateFile(id, newFile)

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
	restaurantId := c.GetString("id")

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
