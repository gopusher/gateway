package log

import (
	"sort"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	// Level is the minimum enabled logging level. Note that this is a dynamic
	// level, so calling Config.Level.SetLevel will atomically change the log
	// level of all loggers descended from this config.
	Level zap.AtomicLevel `json:"level" yaml:"level"`
	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stacktraces more liberally.
	Development bool `json:"development" yaml:"development"`
	// DisableCaller stops annotating logs with the calling function's file
	// name and line number. By default, all logs are annotated.
	// 调用 log方法 的文件名和行号
	DisableCaller bool `json:"disableCaller" yaml:"disableCaller"`
	// DisableStacktrace completely disables automatic stacktrace capturing. By
	// default, stacktraces are captured for WarnLevel and above logs in
	// development and ErrorLevel and above in production.
	DisableStacktrace bool `json:"disableStacktrace" yaml:"disableStacktrace"`
	// Sampling sets a sampling policy. A nil SamplingConfig disables sampling.
	// 日志记录限流
	// 每秒日志数量为Initial，如果超过，则每 Thereafter 条才记录
	// 每秒计数器会重置
	Sampling *zap.SamplingConfig `json:"sampling" yaml:"sampling"`
	// Encoding sets the logger's encoding. Valid values are "json" and
	// "console", as well as any third-party encodings registered via
	// RegisterEncoder.
	Encoder zapcore.Encoder

	WriteSyncer zapcore.WriteSyncer

	InitialFields map[string]interface{} `json:"initialFields" yaml:"initialFields"`
}

func (cfg Config) Build(opts ...zap.Option) *zap.Logger {
	core := zapcore.NewCore(cfg.Encoder, cfg.WriteSyncer, cfg.Level)

	log := zap.New(core, cfg.buildOptions()...)
	if len(opts) > 0 {
		log = log.WithOptions(opts...)
	}

	return log
}

func (cfg Config) buildOptions() []zap.Option {
	opts := make([]zap.Option, 0, 4)

	if cfg.Development {
		opts = append(opts, zap.Development())
	}

	//记录 log 的调用文件名和行号
	if !cfg.DisableCaller {
		opts = append(opts, zap.AddCaller())
		opts = append(opts, zap.AddCallerSkip(1))
	}

	//stacktrace
	//正式环境 error 级别以上记录 stacktrace
	stackLevel := zap.ErrorLevel
	//开发环境 warn 级别以上记录 stacktrace
	if cfg.Development {
		stackLevel = zap.WarnLevel
	}
	if !cfg.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	if cfg.Sampling != nil {
		opts = append(opts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewSampler(core, time.Second, int(cfg.Sampling.Initial), int(cfg.Sampling.Thereafter))
		}))
	}

	if len(cfg.InitialFields) > 0 {
		fs := make([]zap.Field, 0, len(cfg.InitialFields))
		keys := make([]string, 0, len(cfg.InitialFields))
		for k := range cfg.InitialFields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fs = append(fs, zap.Any(k, cfg.InitialFields[k]))
		}
		opts = append(opts, zap.Fields(fs...))
	}

	return opts
}
