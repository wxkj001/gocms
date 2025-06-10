package service

import (
	"gocms/router"
	"gocms/service/internal/admin"
	"gocms/service/internal/install"
	"gocms/service/internal/plugin"
	"gocms/service/internal/test"
	"gocms/service/internal/upload"

	"go.uber.org/fx"
)

var ServiceModule = fx.Module("serviceModule",
	fx.Provide(
		router.AsRoute(admin.NewAdminRouter),
		router.AsRoute(plugin.NewPluginRouter),
		router.AsRoute(upload.NewUploadRouter),
		router.AsRoute(test.NewTestRouter),
		router.AsRoute(install.NewInstallRouter),
	),
)
