package config

import (
	"os"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

var defaultMapStructureDecoderConfig = []viper.DecoderConfigOption{
	func(config *mapstructure.DecoderConfig) {
		config.TagName = "mapstructure"
		//config.TagName = "yaml"
		config.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			func(f reflect.Kind, t reflect.Kind, data interface{}) (interface{}, error) {
				if f != reflect.String || t != reflect.String {
					return data, nil
				}
				return os.ExpandEnv(data.(string)), nil
			},
		)
	},
}

func mapStructureParse(input interface{}, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
	}
	for _, opt := range defaultMapStructureDecoderConfig {
		opt(config)
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	if err := decoder.Decode(input); err != nil {
		return err
	}

	return nil
}

//rawConfig: 未解析得 map
//cfg: 需要解析的 cfg 指针
//tmpInitCfg: 需要解析的 cfg 的实例化，为了进行原子替换的临时指针
func UnmarshalConfig(rawConfig map[string]interface{}, tmpInitCfg interface{}) error {
	// 忽略 .env 是否存在
	// 选用此库是因为 viper 在用，减少更多的库依赖
	_ = gotenv.OverLoad(".env")
	//v.Set("app.redis_prefix_key", v.GetString("app.redis_prefix_key")+":"+v.GetString("symbol")+":")

	if err := mapStructureParse(rawConfig, tmpInitCfg); err != nil {
		return err
	}
	//校验配置
	validate := validator.New()
	if err := validate.Struct(tmpInitCfg); err != nil {
		return err
	}

	return nil
}

func LoadConfig(config interface{}, cfgFile string, flagSet *pflag.FlagSet, keys map[string]string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigFile(cfgFile)
	v.SetConfigType("yaml")
	//命令行参数覆盖
	//配置 key : 命令行参数 name
	for key, name := range keys {
		if err := v.BindPFlag(key, flagSet.Lookup(name)); err != nil {
			return nil, err
		}
	}
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	//watch 配置文件更新
	//go get github.com/fsnotify/fsnotify
	//v.WatchConfig()
	//v.OnConfigChange(func(in fsnotify.Event) {})

	err := UnmarshalConfig(v.AllSettings(), config)
	if err != nil {
		return nil, err
	}

	return v, nil
}
