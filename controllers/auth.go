package controllers

import (
	"fmt"
	"smart-serve/models"
	"smart-serve/utils"

	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type SignInData struct {
	Email    string `json:"email" binding:"required,email" example:"user@gmail.com"`
	Password string `json:"password" binding:"required,min=8" example:"12345678"`
}

type SignInResponse struct {
	AccessToken string `json:"accessToken"`
}

// @Tags Auth
// @Accept json
// @Produce json
// @Param data body SignInData true "Sign in data"
// @Success 200 {object} SignInResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/sign-in [post]
func SignIn(c *gin.Context) {
	var signInData SignInData

	if err := c.ShouldBindJSON(&signInData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := models.GetUserByEmail(signInData.Email)
	fmt.Println(user)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Email not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(signInData.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid password"})
		return
	}

	accessToken, _ := utils.GenerateJWT(user.ID)

	c.JSON(http.StatusOK, SignInResponse{
		AccessToken: accessToken,
	})
}

// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/me [get]
// @Security BearerAuth
func GetMe(c *gin.Context) {
	userId := c.GetString("userId")

	user, err := models.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
