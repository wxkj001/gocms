package app

import (
	"gocms/plugin"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var pluginModule = fx.Module("pluginModule", fx.Provide(NewPlugin))

func NewPlugin(config *viper.Viper) *plugin.Plugins {
	return plugin.New(config)
}
