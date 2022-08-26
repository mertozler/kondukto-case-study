package repository

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/mertozler/internal/config"
)

type Repository struct {
	repo *redis.Client
}

func NewRepository(config *config.Redis) *Repository {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host,
		Password: config.Password,
		DB:       0,
	})
	return &Repository{repo: client}
}

func (r *Repository) SetScanData(key string, value interface{}) error {
	bs, _ := json.Marshal(value)
	err := r.repo.Set(key, bs, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetScanData(key string) (interface{}, error) {
	val, err := r.repo.Get(key).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}
