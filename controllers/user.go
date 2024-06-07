package controllers

import (
	"smart-serve/models"

	"net/http"

	"github.com/gin-gonic/gin"
)

// can use id is string or int

// @Tags Users
// @Accept  json
// @Produce  json
// @Param data body models.CreateUserInput true "User Data"
// @Success 200 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var createUserInput models.CreateUserInput

	if err := c.ShouldBindJSON(&createUserInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user := models.User{
		Name:         createUserInput.Name,
		Email:        createUserInput.Email,
		Password:     createUserInput.Password,
		RestaurantID: createUserInput.RestaurantID,
	}

	user, err := models.CreateUser(user)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Tags Users
// @Accept  json
// @Produce  json
// @Success 200 {array} models.User
// @Router /users [get]
func GetUsers(c *gin.Context) {
	users := models.GetUsers()
	c.JSON(http.StatusOK, users)
}

// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} models.ErrorResponse
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := models.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Param user body models.UpdateUserInput true "User Data"
// @Success 200 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [patch]
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.UpdateUserInput

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	updatedUser, err := models.UpdateUser(id, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := models.DeleteUser(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
