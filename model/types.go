package model

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// Type type
type Type struct {
	cnf configType
	log *zap.SugaredLogger
	rDB *redis.ClusterClient
}

type configType struct {
	Redis redisConfig
}

type redisConfig struct {
	Addrs string
}
