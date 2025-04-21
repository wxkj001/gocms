package test

import (
	"fmt"
	"gocms/middleware"
	"gocms/model"
	"gocms/plugin"
	"gocms/router"
	"gocms/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewTestRouter(log *zap.Logger, config *viper.Viper, plugin *plugin.Plugins, middle middleware.MiddlewareParams, models model.ModelParams) *TestRouter {
	return &TestRouter{log: log, config: config, middle: middle, model: models, plugins: plugin}
}

type TestRouter struct {
	middle  middleware.MiddlewareParams
	log     *zap.Logger
	config  *viper.Viper
	model   model.ModelParams
	g       *gin.Engine
	plugins *plugin.Plugins
}

func (c *TestRouter) RouteRegister(g *gin.Engine, r *gin.RouterGroup) {
	c.g = g
	r.GET("/test2", c.GetCatList2)
	r.GET("/test", c.GetCatList)
}

func (c *TestRouter) GetCatList(ctx *gin.Context) {
	err := c.plugins.Add("reactor")
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code: 500,
			Data: err.Error(),
		})
		return
	}
	pl, err := c.plugins.Use("reactor")
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code: 500,
			Data: err.Error(),
		})
		return
	}
	exit, out, err := pl.Call("greet", []byte("Yellow, World!"))
	if err != nil {
		fmt.Println(err)
		os.Exit(int(exit))
	}
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: string(out),
	})
}
func (c *TestRouter) GetCatList2(ctx *gin.Context) {
	token, _ := utils.GenerateToken(map[string]any{
		"user_id": 1,
		"role":    "1",
	}, c.config)
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: token,
	})
}
