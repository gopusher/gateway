package config

import (
	"gopkg.in/ini.v1"
)

type Config struct {
	value *ini.File
}

func NewConfig(filename string) *Config {
	cfg, err := ini.Load(filename)

	if err != nil {
		panic("加载配置文件失败" + filename + ", 原因: " + err.Error())
	}
	return &Config {
		value: cfg,
	}
}

func (c *Config) Get(key string) (*ini.Key) {
	return c.value.Section("").Key(key)
}
