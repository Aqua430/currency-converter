package cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	if err := Rdb.Ping(Ctx).Err(); err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}
}

func Set(key string, value string, expiration time.Duration) {
	err := Rdb.Set(Ctx, key, value, expiration).Err()
	if err != nil {
		log.Printf("Ошибка при установке ключа %s: %v", key, err)
	}
}

func Get(key string) (string, error) {
	val, err := Rdb.Get(Ctx, key).Result()
	if err == redis.Nil {
		log.Printf("Ключ %s не найден", key)
		return "", nil
	} else if err != nil {
		log.Printf("Ошибка при получении ключа %s: %v", key, err)
	}
	return val, nil
}

func Del(key string) error {
	return Rdb.Del(Ctx, key).Err()
}

func Exists(key string) (bool, error) {
	val, err := Rdb.Exists(Ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val == 1, nil
}
