package app

import (
	"fmt"
	"net/http"

	"github.com/Sarvesh-10/ReadEazeBackend/internal/middleware"
	"github.com/Sarvesh-10/ReadEazeBackend/internal/models"
)

func SSEHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	client := models.CreateClient(userID)
	defer models.RemoveSSEClient(userID)

	for {
		select {
		case msg := <-client.Channel:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
