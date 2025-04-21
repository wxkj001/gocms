package service

import (
	"gocms/router"
	"gocms/service/internal/admin"
	"gocms/service/internal/plugin"
	"gocms/service/internal/test"

	"go.uber.org/fx"
)

var ServiceModule = fx.Module("serviceModule",
	fx.Provide(
		router.AsRoute(admin.NewAdminRouter),
		router.AsRoute(plugin.NewPluginRouter),
		router.AsRoute(test.NewTestRouter),
	),
)
