package Service

import (
	"Tasks/Redis"
	"Tasks/repository"
	"errors"
)

var (
	ErrNotFoundInCache = errors.New("redis: nil")
)

type Service struct {
	repository *repository.Repository
	cache      *Redis.Cache
}

func CreateService(repository *repository.Repository, cache *Redis.Cache) *Service {
	return &Service{
		repository: repository,
		cache:      cache,
	}
}
