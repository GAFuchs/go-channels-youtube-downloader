package main

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type Components struct {
	redisClient *redis.Client
}

func NewComponents() (*Components, error) {
	// Connect to redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping redis")
	}

	return &Components{
		redisClient: client,
	}, nil
}

func (c *Components) Get(key string) (string, error) {
	val, err := c.redisClient.Get(key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}
	return val, nil
}

func (c *Components) GetObject(key string) (interface{}, error) {
	val, err := c.redisClient.Get(key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}
	return val, nil
}

func (c *Components) Set(key string, value string, ttl time.Duration) error {
	err := c.redisClient.Set(key, value, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Components) SetObject(key string, value interface{}, ttl time.Duration) error {
	err := c.redisClient.Set(key, value, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Components) Del(key string) error {
	err := c.redisClient.Del(key).Err()
	if err != nil {
		return err
	}
	return nil
}
