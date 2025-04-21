package plugin

import (
	"gocms/cache"
	"gocms/middleware"
	"gocms/model"
	"gocms/plugin"
	"gocms/router"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewPluginRouter(log *zap.Logger, config *viper.Viper, plugin *plugin.Plugins, middle middleware.MiddlewareParams, models model.ModelParams, e *casbin.Enforcer, cache cache.Cache) *PluginRouter {
	return &PluginRouter{log: log, config: config, middle: middle, model: models, e: e, cache: cache, plugins: plugin}
}

type PluginRouter struct {
	middle  middleware.MiddlewareParams
	log     *zap.Logger
	config  *viper.Viper
	model   model.ModelParams
	g       *gin.Engine
	e       *casbin.Enforcer
	cache   cache.Cache
	plugins *plugin.Plugins
}

// 注册路由
func (c *PluginRouter) RouteRegister(g *gin.Engine, r *gin.RouterGroup) {
	// plugin := r.Group("/plugin")
	r.GET("/admin/plugin/list", c.PluginList)
	r.PUT("/admin/plugin/register/:name", c.PluginRegister)
	r.DELETE("/admin/plugin/register/:name", c.PluginUnregister)
}

// 插件列表
func (c *PluginRouter) PluginList(ctx *gin.Context) {
	ps := []string{}
	filepath.Walk("./public/plugin", func(path string, info os.FileInfo, err error) error {
		c.log.Debug("PluginRouter", zap.String("path", path))
		name := info.Name()
		ss := strings.Split(name, ".")
		if len(ss) > 1 {
			if ss[1] == "wasm" {
				ps = append(ps, ss[0])
			}
		}
		return nil
	})
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: ps,
	})
}

// 注册插件
func (c *PluginRouter) PluginRegister(ctx *gin.Context) {
	name := ctx.Param("name")
	err := c.plugins.Add(name)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	p, err := c.plugins.Use(name)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	_, out, err := p.Call("install", []byte(""))
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: string(out),
	})
}

// 卸载插件
func (c *PluginRouter) PluginUnregister(ctx *gin.Context) {
	name := ctx.Param("name")
	p, err := c.plugins.Use(name)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	_, out, err := p.Call("uninstall", []byte(""))
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	c.plugins.Remove(name)
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: string(out),
	})
}
