package app

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logModule = fx.Module("logModule", fx.Provide(NewLog))

func NewLog(config *viper.Viper) (*zap.Logger, error) {
	fileEncoderConfig := zap.NewProductionEncoderConfig()
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	writer := &lumberjack.Logger{
		Filename:   config.GetString("log.filename"),
		MaxSize:    config.GetInt("log.max_size"), // megabytes
		MaxBackups: config.GetInt("log.max_backups"),
		MaxAge:     config.GetInt("log.max_age"), // days
		Compress:   config.GetBool("log.compress"),
	}
	level, err := zapcore.ParseLevel(config.GetString("log.level"))
	if err != nil {
		return nil, err
	}
	// 创建两个核心，一个用于控制台，一个用于文件
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleEncoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(fileEncoderConfig),
		zapcore.AddSync(writer),
		level,
	)
	// 合并两个核心
	core := zapcore.NewTee(consoleCore, fileCore)
	log := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return log, nil
}
