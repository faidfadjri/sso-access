package database

import (
	"akastra-access/internal/app/config"
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	conf := config.Load()
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Println("Failed to connect to Redis:", err)
		// We might not want to fatal here if Redis is optional, but for a Session Store it's critical
		// log.Fatal("Failed to connect to Redis:", err) 
	} else {
		log.Println("Connected to Redis")
	}

	return rdb
}
