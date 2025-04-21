package app

import (
	"gocms/plugin"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

var pluginModule = fx.Module("pluginModule", fx.Provide(NewPlugin))

func NewPlugin(config *viper.Viper, log *zap.Logger, db *xorm.Engine) *plugin.Plugins {
	return plugin.New(config, log, db)
}
