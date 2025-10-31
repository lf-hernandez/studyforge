package handlers

import (
	"net/http"
	"time"

	"studyforge/pkg/utils"
)

var startTime = time.Now()

// HandleHealth returns health status
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	uptime := int(time.Since(startTime).Seconds())

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"version": "1.0.0-mvp",
		"uptime":  uptime,
	})
}
