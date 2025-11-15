package config

import (
	"context"
	//"fmt"
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
	//REDIS_PORT := viper.GetString("REDIS_PORT")

	if REDIS_HOST == "" {
		REDIS_HOST = "redis://default:YourPassword@your-redis-host.redis.cloud:12345"
	}

	/*opt, err := redis.ParseURL(REDIS_HOST)
	if err != nil {
		log.Fatalf("Invalid Redis URL: %v", err)
	}*/

	//SGku1ECc7k9UWA9OipVduHji3FF1kDPZ

	client := redis.NewClient(&redis.Options{
		Addr:     "redis-10691.c257.us-east-1-3.ec2.cloud.redislabs.com:10691",
		Username: "Indroneel007",
		Password: "SGku1ECc7k9UWA9OipVduHji3FF1kDPZ",
		DB:       0,
	})

	redisCache := &RedisCache{
		Client:           client,
		Store:            persist.NewRedisStore(client),
		DefaultCacheTime: 10 * time.Second,
	}

	err = client.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return redisCache
}

func init() {
	Redis = SetupRedisCache()
	if Redis == nil {
		log.Fatal("Failed to initialize Redis cache")
	} else {
		log.Println("Redis cache initialized successfully")
	}
}
