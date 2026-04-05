package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cfg := loadConfig()

	// Ensure watch directory exists
	if err := os.MkdirAll(cfg.WatchDir, 0755); err != nil {
		log.Fatalf("failed to create watch dir: %v", err)
	}

	db := initDB(cfg.DBPath)
	defer db.Close()

	// Process any existing CSVs on startup
	processExistingCSVs(cfg.WatchDir, db, cfg.NtfyTopic)

	// Start file watcher in background
	go watchDirectory(cfg.WatchDir, db, cfg.NtfyTopic)

	r := gin.Default()
	r.GET("/health", healthHandler())
	r.GET("/api/measurements", getMeasurements(db))

	log.Printf("starting server on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
