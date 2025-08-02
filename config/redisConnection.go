package config

import (
	"context"
	"fmt"
	"log"

	//"os"
	"time"

	"github.com/chenyahui/gin-cache/persist"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	//"github.com/subosito/gotenv"
)

type RedisCache struct {
	Client           *redis.Client
	Store            *persist.RedisStore
	DefaultCacheTime time.Duration
}

var ctx = context.Background()

var Redis *RedisCache

func SetupRedisCache() *RedisCache {
	/*err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}*/
	paths := []string{".env", "../.env", "../../.env"}
	for _, path := range paths {
		_ = godotenv.Load(path)
	}

	var err error

	viper.AutomaticEnv()

	REDIS_HOST := viper.GetString("REDIS_HOST")
	REDIS_PORT := viper.GetString("REDIS_PORT")

	if REDIS_HOST == "" {
		REDIS_HOST = "red-d2763u63jp1c73edni1g"
	}

	if REDIS_PORT == "" {
		REDIS_PORT = "6379"
	}

	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr: fmt.Sprintf(
			"%s:%s",
			REDIS_HOST,
			REDIS_PORT,
		),
	})

	redis := &RedisCache{
		Client:           client,
		Store:            persist.NewRedisStore(client),
		DefaultCacheTime: 10 * time.Second,
	}

	err = client.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return redis
}

func init() {
	Redis = SetupRedisCache()
	if Redis == nil {
		log.Fatal("Failed to initialize Redis cache")
	} else {
		log.Println("Redis cache initialized successfully")
	}
}
