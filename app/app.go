package app

import (
	"gocms/extend"
	"gocms/middleware"
	"gocms/model"
	"gocms/service"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func New() *fx.App {
	return fx.New(
		logModule,
		// 注入config模块
		configModule,
		// 注入db模块
		dbModule,
		// 注入http模块
		httpModule,
		// 注入缓存模块
		cacheModule,
		// 注入业务模块
		service.ServiceModule,
		// 注入model模块
		model.ModelModule,
		// 注入中间件模块
		middleware.MiddlewareModule,
		// 注入casbin
		casbinModule,
		// 注入extend模块
		extend.ExtendModule,
		extend.ResponseModules,
		//
		// validates.ValidateModule,
		//	定时模块
		cronModule,
		// 注入插件模块
		pluginModule,
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	)
}
