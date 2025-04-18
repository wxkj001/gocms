package app

import (
	"context"
	"crypto/tls"
	"gocms/middleware"
	"gocms/router"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var httpModule = fx.Module("httpServer", fx.Provide(fx.Annotate(NewHttp, fx.ParamTags(`group:"routes"`))), fx.Invoke(func(*gin.Engine) {}))

func NewHttp(routers []router.Route, config *viper.Viper, lc fx.Lifecycle, log *zap.Logger) *gin.Engine {
	web := gin.Default()
	web.Use(middleware.NewLogger(config, log).Logger)
	web.Use(gin.Recovery())
	/* web.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, x-os",
	})) */
	api := web.Group("/api", middleware.NewDataEncrypt(config).Encrypt)
	for _, v := range routers {
		v.RouteRegister(web, api)
	}
	var tlsc *tls.Config
	if config.GetBool("web.tls.enable") {
		tlsc = &tls.Config{
			ServerName: config.GetString("web.tls.servername"),
			MinVersion: tls.VersionTLS11,
			MaxVersion: tls.VersionTLS13,
		}
	}
	srv := &http.Server{
		Addr:      ":" + config.GetString("web.port"),
		Handler:   web,
		TLSConfig: tlsc,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if config.GetBool("web.tls.enenable") {
				go srv.ListenAndServeTLS(config.GetString("web.tls.cert"), config.GetString("web.tls.key"))
			} else {
				go srv.ListenAndServe()
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return web
}
