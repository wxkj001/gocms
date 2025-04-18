package app

import (
	"context"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var configModule = fx.Module("configModule", fx.Provide(NewConfig))

func NewConfig(lc fx.Lifecycle) (*viper.Viper, error) {
	config := viper.New()
	config.SetConfigName("config")        // name of config file (without extension)
	config.SetConfigType("yaml")          // REQUIRED if the config file does not have the extension in the name
	config.AddConfigPath("/etc/gocms/")   // path to look for the config file in
	config.AddConfigPath("$HOME/.gocms/") // call multiple times to add many search paths
	config.AddConfigPath("./config")
	config.AddConfigPath(".")    // optionally look for config in the working directory
	err := config.ReadInConfig() // Find and read the config file
	if err != nil {              // Handle errors reading the config file
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			config.OnConfigChange(func(in fsnotify.Event) {
				err := config.ReadInConfig() // Find and read the config file
				if err != nil {              // Handle errors reading the config file
					return
				}
			})
			return nil
		},
	})
	viper.SetOptions()
	return config, nil
}
