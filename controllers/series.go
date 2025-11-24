package controllers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/OR-Sasaki/cat-backend/models"
)

func SeriesRouter(router *gin.RouterGroup) {
	series := router.Group("/series")
	{
		series.GET("", GetAllSeries)
	}
}

// **************************************************
// GetAllSeries
// **************************************************

type SeriesResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func GetAllSeries(c *gin.Context) {
	series, err := models.GetAllSeries(c.Request.Context())
	if err != nil {
		slog.Error("failed to get all series", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	response := make([]SeriesResponse, len(series))
	for i, s := range series {
		response[i] = SeriesResponse{
			ID:   s.ID,
			Name: s.Name,
		}
	}

	c.JSON(http.StatusOK, response)
}
