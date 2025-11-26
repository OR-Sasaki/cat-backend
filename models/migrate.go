package models

import "github.com/OR-Sasaki/cat-backend/config"

func Migrate() error {
	return config.DB.AutoMigrate(&User{}, &Series{}, &Outfit{})
}
