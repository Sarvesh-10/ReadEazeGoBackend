package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"context"

	"github.com/Sarvesh-10/ReadEazeBackend/internal/models"
	"github.com/Sarvesh-10/ReadEazeBackend/utility"
	"github.com/redis/go-redis/v9"
)

type CacheRepository interface {
	SaveUserBookProfile(ctx context.Context, userID int, bookID int, profileData models.UserBookProfile) error
	GetUserBookProfile(userID int, bookID int) (string, error)
	DeleteUserBookProfile(userID int, bookID int) error
	PushToQueue(queueName string, data []byte) error
	LeftPush(queueName string, bookIndexingJob models.BookIndexingJob) error
}

type RedisBookCache struct {
	redis  *redis.Client
	logger *utility.Logger
}

func NewRedisBookCache(redisClient *redis.Client, logger *utility.Logger) *RedisBookCache {
	return &RedisBookCache{
		redis:  redisClient,
		logger: logger,
	}
}
func (r *RedisBookCache) PushToQueue(queueName string, data []byte) error {
	return r.redis.RPush(context.Background(), queueName, data).Err()
}
func (r *RedisBookCache) LeftPush(queueName string, bookIndexingJob models.BookIndexingJob) error {
	data, err := json.Marshal(bookIndexingJob)
	if err != nil {
		return err
	}
	return r.redis.LPush(context.Background(), queueName, data).Err()
}
func (r *RedisBookCache) SaveUserBookProfile(ctx context.Context, userID int, bookID int, profileData models.UserBookProfile) error {
	key := r.getUserBookProfileKey(userID, bookID)
	data, err := profileData.ToJSON()
	if err != nil {
		return err
	}
	return r.redis.Set(ctx, key, data, 2*time.Hour).Err()
}

func (r *RedisBookCache) GetUserBookProfile(userID int, bookID int) (string, error) {
	key := r.getUserBookProfileKey(userID, bookID)
	data, err := r.redis.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return data, nil
}
func (r *RedisBookCache) getUserBookProfileKey(userID int, bookID int) string {
	return fmt.Sprintf("user:%d:book:%d:profile", userID, bookID)
}

func (r *RedisBookCache) DeleteUserBookProfile(userID int, bookID int) error {
	key := r.getUserBookProfileKey(userID, bookID)
	return r.redis.Del(context.Background(), key).Err()
}
