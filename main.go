package main

import (
	"log/slog"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/OR-Sasaki/cat-backend/config"
	"github.com/OR-Sasaki/cat-backend/models"
	"github.com/OR-Sasaki/cat-backend/routers"
)

func main() {

	var err error
	if config.DB, err = gorm.Open(sqlite.Open(config.DBPath), &gorm.Config{}); err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	if err = models.Migrate(); err != nil {
		slog.Error("failed to run migration", "error", err)
		os.Exit(1)
	}

	router := routers.SetupRouter()
	if err = router.Run(":" + config.Port); err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
