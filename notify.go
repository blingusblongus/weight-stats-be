package main

import (
	"log"
	"net/http"
	"strings"
)

func sendNotification(topic, message string) {
	if topic == "" {
		return
	}

	req, err := http.NewRequest("POST", topic, strings.NewReader(message))
	if err != nil {
		log.Printf("ntfy: failed to create request: %v", err)
		return
	}
	req.Header.Set("Title", "Weight Stats")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("ntfy: failed to send: %v", err)
		return
	}
	resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Printf("ntfy: unexpected status %d", resp.StatusCode)
	}
}
