package services

import (
	"webinar/graphql/server/internal/services/storage"
	"webinar/graphql/server/internal/system/database/redis"
)

// DI for services
type Container struct {
	SessionStorage storage.Session
}

func Setup(redisClient redis.Repository) *Container {
	sc := &Container{}
	sc.SessionStorage = storage.NewSessionStorage(redisClient)
	return sc
}
