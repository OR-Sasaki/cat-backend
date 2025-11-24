package models

import (
	"context"

	"gorm.io/gorm"

	"github.com/OR-Sasaki/cat-backend/config"
)

type OutfitType string

const (
	OutfitTypeDefault OutfitType = "default"
)

type Outfit struct {
	gorm.Model
	Name      string     `gorm:"not null"`
	Type      OutfitType `gorm:"not null"`
	SeriesID  uint       `gorm:"not null"`
	AssetPath string
	Series    Series
}

func GetOutfit(ctx context.Context, id uint) (*Outfit, error) {
	outfit, err := gorm.G[Outfit](config.DB).Preload("Series", nil).Where("id = ?", id).First(ctx)
	return &outfit, err
}

func GetAllOutfits(ctx context.Context, seriesID *uint) ([]Outfit, error) {
	query := gorm.G[Outfit](config.DB)
	if seriesID != nil {
		outfits, err := query.Where("series_id = ?", *seriesID).Find(ctx)
		return outfits, err
	}
	outfits, err := query.Find(ctx)
	return outfits, err
}
