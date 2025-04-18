package install

import (
	"gocms/middleware"
	"gocms/model"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewInstallRouter(log *zap.Logger, config *viper.Viper, middle middleware.MiddlewareParams, models model.ModelParams) *InstallRouter {
	return &InstallRouter{log: log, config: config, middle: middle, model: models}
}

type InstallRouter struct {
	middle middleware.MiddlewareParams
	log    *zap.Logger
	config *viper.Viper
	model  model.ModelParams
}

func (c *InstallRouter) RouteRegister(g *gin.Engine, r *gin.RouterGroup) {
	r.GET("/install", c.Install)
}
func (c *InstallRouter) Install(ctx *gin.Context) {
	c.log.Info("install")
}
