package config

import "github.com/gopusher/gateway/pkg/log"

type AppConfig struct {
	AppName      string            `mapstructure:"app_name"`
	AppDebug     bool              `mapstructure:"app_debug"`
	LoggerConfig *log.LoggerConfig `mapstructure:"logging" validate:"required,dive"`
}

func (appCfg *AppConfig) InitLoggerConfig() {
	appCfg.LoggerConfig.Development = appCfg.AppDebug
	appCfg.LoggerConfig.AppName = appCfg.AppName
}
