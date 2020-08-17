package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Repository struct {
	*redis.Client
}

func NewLocal() *Repository {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := c.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return &Repository{c}
}
