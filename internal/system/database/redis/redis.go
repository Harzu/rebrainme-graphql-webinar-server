package redis

import (
	"github.com/go-redis/redis/v8"
)

type Repository interface {
	GetConnection() *redis.ClusterClient
}

type repository struct {
	client *redis.ClusterClient
}

func New(config *Config) (Repository, error) {
	r := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{config.URL},
		Password: config.Pass,
	})
	return &repository{client: r}, nil
}

func (r *repository) GetConnection() *redis.ClusterClient {
	return r.client
}
