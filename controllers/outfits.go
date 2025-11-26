package controllers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/OR-Sasaki/cat-backend/authenticate"
	"github.com/OR-Sasaki/cat-backend/models"
)

func OutfitsRouter(router *gin.RouterGroup) {
	outfits := router.Group("/outfits")
	{
		authenticate.GETWithAuth(outfits, "", GetAllOutfits)
	}
}

// **************************************************
// GetAllOutfits
// **************************************************

type OutfitResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	SeriesID  uint   `json:"series_id"`
	AssetPath string `json:"asset_path"`
}

func GetAllOutfits(c *gin.Context, u *models.User) {
	outfits, err := models.GetAllOutfits(c.Request.Context())
	if err != nil {
		slog.Error("failed to get all outfits", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	response := make([]OutfitResponse, len(outfits))
	for i, outfit := range outfits {
		response[i] = OutfitResponse{
			ID:        outfit.ID,
			Name:      outfit.Name,
			Type:      string(outfit.Type),
			SeriesID:  outfit.SeriesID,
			AssetPath: outfit.AssetPath,
		}
	}

	c.JSON(http.StatusOK, response)
}
