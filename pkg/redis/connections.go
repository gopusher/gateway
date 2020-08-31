package redis

import (
	"github.com/gopusher/gateway/pkg/log"
	"fmt"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type Configs map[string]*Config

type Connections struct {
	configs     Configs
	connections map[string]*redis.Client
}

func InitConnections(configs Configs) *Connections {
	connections := make(map[string]*redis.Client)
	for conn, config := range configs {
		redisPool, err := NewRedis(config)
		if err != nil {
			log.Panic(fmt.Sprintf("newRedis redis, conn: %s, error: %s", conn, err.Error()), zap.Error(err))
		}

		connections[conn] = redisPool
	}

	return &Connections{
		configs:     configs,
		connections: connections,
	}
}

func (conns *Connections) Connection(conn string) *redis.Client {
	if conn == "" {
		conn = "default"
	}

	return conns.connections[conn]
}

func (conns *Connections) Publish(conn string, channel string, message interface{}) error {
	if err := conns.Connection(conn).Publish(channel, message).Err(); err != nil {
		log.Error("redis publish error, channel: "+channel, zap.Error(err), zap.Any("message", message))
		return err
	}

	log.Debug("redis publish channel:"+channel, zap.Any("message", message))
	return nil
}
