package cfg

import (
	"github.com/gopusher/gateway/app/gateway/app/api"
	"github.com/gopusher/gateway/app/gateway/app/protocols"
	"github.com/gopusher/gateway/pkg/config"
	"github.com/gopusher/gateway/pkg/dingtalk"
	"github.com/gopusher/gateway/pkg/redis"
)

var Config = &AppConf{}

type AppConf struct {
	config.AppConfig `mapstructure:",squash"`

	ApiServer *api.Config `mapstructure:"api_server" validate:"required"`

	DingTalk dingtalk.Configs `mapstructure:"dingtalk"`

	Redis redis.Configs `mapstructure:"redis" validate:"dive"`

	Node string `mapstructure:"node" validate:"required"`

	Server map[string]map[string]interface{} `mapstructure:"server" validate:"required,len=1"`
}

func (cfg *AppConf) Protocol() string {
	for name := range cfg.Server {
		return name
	}

	panic("undefined market name")
}

func (cfg *AppConf) Unpack(server protocols.Server) error {
	return config.UnmarshalConfig(cfg.Server[server.Protocol()], server.Config())
}
