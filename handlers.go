package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func healthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func getMeasurements(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := c.Query("start")
		end := c.Query("end")

		measurements, err := queryMeasurements(db, start, end)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if measurements == nil {
			measurements = []Measurement{}
		}
		c.JSON(http.StatusOK, measurements)
	}
}
