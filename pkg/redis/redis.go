package redis

import (
	"errors"

	"github.com/go-redis/redis"
)

type Config struct {
	Addr     string `mapstructure:"addr" validate:"required"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
}

func NewRedis(config *Config) (*redis.Client, error) {
	redisPool := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password, // no password set
		DB:       config.Db,       // use default DB
		//DialTimeout:  10 * time.Second,
		//ReadTimeout:  30 * time.Second,
		//WriteTimeout: 30 * time.Second,
		//PoolSize:     10,
		//PoolTimeout:  30 * time.Second,
	})

	if err := redisPool.Ping().Err(); err != nil {
		return nil, errors.New("redis 连接失败: " + err.Error())
	}

	return redisPool, nil
}
