package websocket

type Config struct {
	AppKey    string `mapstructure:"app_key" validate:"required"`
	AppSecret string `mapstructure:"app_secret" validate:"required"`

	Address     string `mapstructure:"address" validate:"required"`
	Ssl         bool   `mapstructure:"ssl"`
	SslCertFile string `mapstructure:"ssl_cert_file" validate:"required_with=Ssl"`
	SslKeyFile  string `mapstructure:"ssl_key_file" validate:"required_with=Ssl"`

	ClientIdAlias string `mapstructure:"client_id_alias" validate:"required"`
	TokenAlias    string `mapstructure:"token_alias" validate:"required"`
	TimeAlias     string `mapstructure:"time_alias" validate:"required"`
	TimeWindow    int64  `mapstructure:"time_window" validate:"required"`
}

var defaultConfig = &Config{}
