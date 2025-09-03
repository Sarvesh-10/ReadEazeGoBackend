package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Sarvesh-10/ReadEazeBackend/internal/domain"
	"github.com/Sarvesh-10/ReadEazeBackend/internal/models"
	"github.com/redis/go-redis/v9"
)

// OUTPUT_JOB_QUEUE = "output_jobs_queue"
func ListenStatusQueue(cache *redis.Client, statusService *StatusService) {
	for {
		result, err := cache.BLPop(context.Background(), 0, "output_jobs_queue").Result()
		if err != nil {
			log.Println("Error reading from queue:", err)
			continue
		}

		var job models.BookIndexingJob
		if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
			log.Println("Invalid message format:", err)
			continue
		}

		client, exists := models.GetClient(job.UserID)
		if !exists {
			log.Println("No active SSE client for user:", job.UserID)
			continue
		}

		message := statusService.ProcessJob(job)
		message.Message = GenerateJobMessage(message)

		select {
		case client.Channel <- message.Message:
			log.Printf("Sent update to user %d for book %d\n", job.UserID, job.BookID)
		default:
			log.Println("Client channel full, skipping update for user:", job.UserID)
		}
	}
}

func GenerateJobMessage(msg domain.SSEMessage) string {
	status := strings.ToLower(string(msg.Status))

	if msg.Name != "" {
		return fmt.Sprintf("Book Indexing %s for %s", status, msg.Name)
	}

	return fmt.Sprintf("Book Indexing %s", status)
}
