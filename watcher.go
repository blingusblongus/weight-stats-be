package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func processCSVFile(path string, db *sql.DB) int64 {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("watcher: failed to open %s: %v", path, err)
		return 0
	}
	defer f.Close()

	measurements, err := parseCSV(f)
	if err != nil {
		log.Printf("watcher: failed to parse %s: %v", path, err)
		return 0
	}

	inserted, err := insertMeasurements(db, measurements)
	if err != nil {
		log.Printf("watcher: failed to insert from %s: %v", path, err)
		return 0
	}

	filename := filepath.Base(path)
	log.Printf("watcher: %s — %d rows parsed, %d new", filename, len(measurements), inserted)
	return inserted
}

func processExistingCSVs(dir string, db *sql.DB, ntfyTopic string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("watcher: failed to read directory %s: %v", dir, err)
		return
	}

	var total int64
	for _, entry := range entries {
		if entry.IsDir() || !isCSV(entry.Name()) {
			continue
		}
		total += processCSVFile(filepath.Join(dir, entry.Name()), db)
	}

	if total > 0 {
		sendNotification(ntfyTopic, fmt.Sprintf("Processed %d new measurement(s) from %d file(s)", total, len(entries)))
	}
}

func watchDirectory(dir string, db *sql.DB, ntfyTopic string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("watcher: failed to create: %v", err)
	}
	defer watcher.Close()

	if err := watcher.Add(dir); err != nil {
		log.Fatalf("watcher: failed to watch %s: %v", dir, err)
	}

	log.Printf("watcher: watching %s for CSV files", dir)

	// Debounce: track last event time per file
	pending := make(map[string]time.Time)
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if !isCSV(event.Name) {
				continue
			}
			if event.Op&(fsnotify.Create|fsnotify.Write) != 0 {
				pending[event.Name] = time.Now()
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("watcher: error: %v", err)

		case <-ticker.C:
			now := time.Now()
			var totalInserted int64
			var processed int
			for path, lastEvent := range pending {
				if now.Sub(lastEvent) >= time.Minute {
					totalInserted += processCSVFile(path, db)
					delete(pending, path)
					processed++
				}
			}
			if totalInserted > 0 {
				sendNotification(ntfyTopic, fmt.Sprintf("Processed %d new measurement(s) from %d file(s)", totalInserted, processed))
			}
		}
	}
}

func isCSV(name string) bool {
	return strings.HasSuffix(strings.ToLower(name), ".csv")
}
