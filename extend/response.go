package extend

import (
	"gocms/cache"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ResponseModule struct {
	Log    *zap.Logger
	Config *viper.Viper
	G      *gin.Engine
	Casbin *casbin.Enforcer
	Cache  cache.Cache
}

// ResponseModules模块
var ResponseModules = fx.Module("responseModule", fx.Provide(NewResponseModule))

type ResponseResult struct {
	fx.Out
	Response *ResponseModule
}
type ResponseParams struct {
	fx.In
	*ResponseModule
}

func NewResponseModule(log *zap.Logger, config *viper.Viper, e *casbin.Enforcer, cache cache.Cache) (ResponseResult, error) {
	return ResponseResult{Response: &ResponseModule{
		Log:    log,
		Config: config,
		G:      gin.Default(),
		Casbin: e,
		Cache:  cache,
	}}, nil
}
