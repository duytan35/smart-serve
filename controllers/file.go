package controllers

import (
	"net/http"
	"smart-serve/models"
	"smart-serve/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 200 {object} Response{data=models.File}
// @Router /files [post]
// @Security BearerAuth
func Upload(c *gin.Context) {
	restaurantId := c.GetString("id")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
		return
	}

	newFile, err := utils.Uploader.UploadFile(file, header)
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
		// TODO: Delete file from storage
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    fileResponse,
	})
}
