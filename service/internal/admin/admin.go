package admin

import (
	"gocms/cache"
	"gocms/middleware"
	"gocms/model"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewAdminRouter(log *zap.Logger, config *viper.Viper, middle middleware.MiddlewareParams, models model.ModelParams, e *casbin.Enforcer, cache cache.Cache) *AdminRouter {
	return &AdminRouter{log: log, config: config, middle: middle, model: models, e: e, cache: cache}
}

type AdminRouter struct {
	middle middleware.MiddlewareParams
	log    *zap.Logger
	config *viper.Viper
	model  model.ModelParams
	g      *gin.Engine
	e      *casbin.Enforcer
	cache  cache.Cache
}

// 注册路由
func (c *AdminRouter) RouteRegister(g *gin.Engine, r *gin.RouterGroup) {
	c.g = g
	admin := r.Group("/admin")
	// 用户
	user := NewUser(c)
	admin.POST("/login", user.Login)
	admin.GET("/captcha", user.Captcha)
	{
		userRouter := admin.Group("/user", c.middle.Jwt.AdminJWT())
		userRouter.GET("/detail", user.Detail)
		userRouter.GET("/list", user.List)
		userRouter.POST("/", user.Add)
		userRouter.POST("/update", user.Update)
		userRouter.DELETE("/:id", user.Delete)
		userRouter.PATCH("/password/reset/:id", user.ResetPassword)
		userRouter.GET("/permissions", user.Permission)
		userRouter.GET("/refresh/token", user.RefreshToken)
	}
	// 权限设置
	permission := NewPermission(c)
	{
		permissionRouter := admin.Group("/permission", c.middle.Jwt.AdminJWT())
		permissionRouter.GET("/tree", permission.Tree)
		permissionRouter.GET("/:type/:id", permission.GetTypeList)
		permissionRouter.GET("/list", permission.List)
		permissionRouter.POST("/", permission.Add)
		permissionRouter.PATCH("/:id", permission.Update)
		permissionRouter.DELETE("/:id", permission.Delete)
		permissionRouter.GET("/apis", permission.Routers)
	}
	// 角色
	role := NewRole(c)
	{
		roleRouter := admin.Group("/role", c.middle.Jwt.AdminJWT())
		roleRouter.GET("/tree", role.Tree)
		roleRouter.GET("/list", role.List)
		roleRouter.POST("/", role.Add)
		roleRouter.PATCH("/:id", role.Update)
		roleRouter.DELETE("/:id", role.Delete)
	}
	// 配置
	config := NewSysConfig(c)
	{
		configRouter := admin.Group("/config", c.middle.Jwt.AdminJWT())
		configRouter.GET("/list", config.List)
		configRouter.PATCH("/:key", config.Update)
		/* configRouter.POST("/", config.Add)
		configRouter.PATCH("/:id", config.Update)
		configRouter.DELETE("/:id", config.Delete) */
	}
}
