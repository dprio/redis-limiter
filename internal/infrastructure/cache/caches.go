package cache

import (
	"github.com/dprio/redis-limiter/internal/infrastructure/config"
)

type Caches struct {
	Client Client
}

func NewCaches(cfg *config.Config) *Caches {

	return &Caches{
		Client: NewRedisClient(cfg.Redis),
	}
}
