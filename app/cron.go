package app

import (
	"context"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var cronModule = fx.Module("cronModule", fx.Provide(newCron), fx.Invoke(func(*cron.Cron) {}))

func newCron(config *viper.Viper, lc fx.Lifecycle, log *zap.Logger) *cron.Cron {
	c := cron.New()
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("cron start")
			c.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("cron stop")
			c.Stop()
			return nil
		},
	})
	return c
}
