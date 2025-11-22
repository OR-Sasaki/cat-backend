package controllers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/OR-Sasaki/cat-backend/models"
)

type UserRegisterRequest struct {
	Name string `json:"name" binding:"required,min=4,max=20"`
}

type UserRegisterResponse struct {
	ID       uint   `json:"id"`
	Password string `json:"password"`
}

func UserRegister(c *gin.Context) {
	var request UserRegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		slog.Error("failed to bind request parameters", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "parameter error"})
		return
	}

	user, password, err := models.RegisterUser(c.Request.Context(), request.Name)
	if err != nil {
		slog.Error("failed to register user", "error", err, "name", request.Name)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, UserRegisterResponse{
		ID:       user.ID,
		Password: password,
	})
}
