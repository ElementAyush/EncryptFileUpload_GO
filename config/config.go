package config

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/go-redis/redis"
	"github.com/minio/minio-go"
)

var minioClient *minio.Client
var client *redis.Client
var minionce sync.Once
var redisonce sync.Once

/**
* This function initialize the minio client and
* Implements singleton pattern
* @return *minio.Client
**/
func Minioconfig() *minio.Client {

	minionce.Do(func() {
		endpoint := os.Getenv("BUCKET_ENDPOINT")
		accessKeyID := os.Getenv("BUCKET_ACCESS_KEY")
		secretAccessKey := os.Getenv("BUCKET_SECRET_KEY")
		useSSL, _ := strconv.ParseBool(os.Getenv("BUCKET_USESSL"))

		minio, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
		minioClient = minio
		if err != nil {
			log.Println("Erro", err)

		}
	})
	return minioClient

}

func RedisClient() *redis.Client {
	redisonce.Do(func() {
		log.Println("Initializing redis client")
		addr := os.Getenv("REDIS_SERVER")
		pass := os.Getenv("REDIS_PASSWORD")
		redisclient := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pass,
			DB:       0,
		})
		log.Println("Initialization successfull")
		client = redisclient
	})
	return client
}
