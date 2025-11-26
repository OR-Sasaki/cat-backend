package models

import (
	"context"

	"gorm.io/gorm"

	"github.com/OR-Sasaki/cat-backend/config"
)

type Series struct {
	gorm.Model
	Name string `gorm:"not null"`
}

func GetAllSeries(ctx context.Context) ([]Series, error) {
	series, err := gorm.G[Series](config.DB).Find(ctx)
	return series, err
}
