package uniqid

import (
	"github.com/gopusher/gateway/pkg/log"

	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

var node *snowflake.Node

func init() {
	var err error
	// Create a new Node with a Node number of 1
	node, err = snowflake.NewNode(1)
	if err != nil {
		log.Panic("snowflake NewNode error", zap.Error(err))
	}
}

func SnowflakeId() string {
	return node.Generate().Base32()
}
