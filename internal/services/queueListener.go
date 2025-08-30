package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Sarvesh-10/ReadEazeBackend/internal/models"
	"github.com/redis/go-redis/v9"
)

// OUTPUT_JOB_QUEUE = "output_jobs_queue"
func ListenStatusQueue(cache *redis.Client) {
	for {
		log.Println("Listening to queue...")
		result, err := cache.BLPop(context.Background(), 0, "output_jobs_queue").Result()
		if err != nil {
			log.Println("Error reading from queue:", err)
			continue
		}

		var statusMsg models.BookIndexingJob

		err = json.Unmarshal([]byte(result[1]), &statusMsg)
		if err != nil {
			log.Println("Invalid message format:", err)
			continue
		}
		log.Println("Received message for user:", statusMsg.UserID)
		client, exists := models.GetClient(statusMsg.UserID)

		if exists {
			msgBytes, _ := json.Marshal(statusMsg)
			select {
			case client.Channel <- string(msgBytes):
			default:
				log.Println("Client channel full, skipping update for", statusMsg.UserID)
			}
		} else {
			// jsonMsg, _ := json.Marshal(statusMsg)
			// cache.RPush(context.Background(), "output_jobs_queue", jsonMsg)
			log.Println("No active client for user:", statusMsg.UserID)
		}
	}
}
