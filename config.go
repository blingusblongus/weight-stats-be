package main

import (
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Port      string
	DBPath    string
	WatchDir  string
	NtfyTopic string
}

func loadConfig() Config {
	cfg := Config{
		Port:      os.Getenv("PORT"),
		DBPath:    os.Getenv("DB_PATH"),
		WatchDir:  os.Getenv("WATCH_DIR"),
		NtfyTopic: os.Getenv("NTFY_TOPIC"),
	}
	if cfg.Port == "" {
		cfg.Port = "8083"
	}
	if cfg.DBPath == "" {
		cfg.DBPath = "./weight-stats.db"
	}
	if cfg.WatchDir == "" {
		cfg.WatchDir = "./watch"
	}
	cfg.WatchDir = expandHome(cfg.WatchDir)
	cfg.DBPath = expandHome(cfg.DBPath)
	return cfg
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}
