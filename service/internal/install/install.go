package install

import (
	"crypto/md5"
	"fmt"
	"gocms/middleware"
	"gocms/model"
	"gocms/router"
	"gocms/utils"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewInstallRouter(log *zap.Logger, config *viper.Viper, middle middleware.MiddlewareParams, models model.ModelParams, casbin *casbin.Enforcer) *InstallRouter {
	return &InstallRouter{log: log, config: config, middle: middle, model: models, Casbin: casbin}
}

type InstallRouter struct {
	middle middleware.MiddlewareParams
	log    *zap.Logger
	config *viper.Viper
	Casbin *casbin.Enforcer
	model  model.ModelParams
}

func (c *InstallRouter) RouteRegister(g *gin.Engine, r *gin.RouterGroup) {
	r.POST("/install", c.Install)
	r.GET("/install/:key", c.Install)
}
func (c *InstallRouter) Install(ctx *gin.Context) {
	c.log.Info("install")
	type install struct {
		Type string           `json:"type"`
		User *model.UserGroup `json:"user"`
	}
	idata := install{}
	if err := ctx.ShouldBind(&idata); err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	switch idata.Type {
	case "adduser":
		err := c.AddPermission(ctx)
		if err != nil {
			ctx.JSON(200, router.Response{
				Code:    500,
				Message: err.Error(),
			})
			return
		}
		err = c.AddUser(ctx, idata.User)
		if err != nil {
			ctx.JSON(200, router.Response{
				Code:    500,
				Message: err.Error(),
			})
			return
		}

	}
	ctx.JSON(200, router.Response{
		Code:    200,
		Message: "success",
	})
}
func (c *InstallRouter) AddUser(ctx *gin.Context, userReq *model.UserGroup) error {

	userReq.Createtime = time.Now()
	userReq.Updatetime = time.Now()
	userReq.Password = fmt.Sprintf("%x", md5.Sum([]byte(userReq.Password)))
	err := c.model.UserModel.CreateUser(userReq)
	if err != nil {
		return err
	}
	role := &model.Role{
		Code: "admin",
		Name: "超级管理员",
	}
	_, err = c.model.RoleModel.CreateRole(role)
	if err != nil {
		return err
	}
	c.Casbin.AddPolicy(utils.ToString(role.ID), "/api/admin/*", "GET|POST|PUT|DELETE|PATCH")

	c.Casbin.AddGroupingPolicy(utils.ToString(userReq.User.ID), utils.ToString(role.ID), "user")
	c.Casbin.LoadPolicy()
	c.Casbin.SavePolicy()
	return nil
}
func (c *InstallRouter) AddPermission(ctx *gin.Context) error {
	ms := []model.Permission{
		{
			ID:        1,
			Parentid:  0,
			Name:      "系统管理",
			Code:      "SysMgt",
			Type:      "MENU",
			Path:      "",
			Icon:      "i-fe:grid",
			Component: "",
			Show:      true,
			Enable:    true,
			Order:     1,
		}, {
			ID:        2,
			Parentid:  0,
			Name:      "个人资料",
			Code:      "UserProfile",
			Type:      "MENU",
			Path:      "/profile",
			Icon:      "i-fe:user",
			Component: "/src/views/profile/index.vue",
			Show:      false,
			Enable:    true,
			Order:     1,
		}, {
			ID:        3,
			Parentid:  0,
			Name:      "基础功能",
			Code:      "Base",
			Type:      "MENU",
			Path:      "/base",
			Icon:      "i-fe:grid",
			Component: "",
			Show:      true,
			Enable:    true,
			Order:     1,
		}, {
			ID:        4,
			Parentid:  1,
			Name:      "资源管理",
			Code:      "Resource_Mgt",
			Type:      "MENU",
			Path:      "/pms/resource",
			Icon:      "i-fe:list",
			Component: "/src/views/pms/resource/index.vue",
			Show:      true,
			Enable:    true,
			Order:     1,
		}, {
			ID:        5,
			Parentid:  1,
			Name:      "角色管理",
			Code:      "RoleMgt",
			Type:      "MENU",
			Path:      "/pms/role",
			Icon:      "i-fe:user-check",
			Component: "/src/views/pms/role/index.vue",
			Show:      true,
			Enable:    true,
			Order:     1,
		}, {
			ID:        6,
			Parentid:  1,
			Name:      "用户管理",
			Code:      "UserMgt",
			Type:      "MENU",
			Path:      "/pms/user",
			Icon:      "i-fe:user",
			Component: "/src/views/pms/user/index.vue",
			Keepalive: true,
			Show:      true,
			Enable:    true,
			Order:     1,
		}, {
			ID:        7,
			Parentid:  1,
			Name:      "数据对象",
			Code:      "UDO",
			Type:      "MENU",
			Path:      "/pm/udo",
			Icon:      "i-fe:server",
			Component: "/src/views/pms/udo/index.vue",
			Keepalive: false,
			Show:      true,
			Enable:    true,
			Order:     1,
		}, {
			ID:        8,
			Parentid:  5,
			Name:      "创建新用户",
			Code:      "AddUser",
			Type:      "BUTTON",
			Path:      "",
			Icon:      "",
			Component: "",
			Keepalive: false,
			Show:      true,
			Enable:    true,
			Order:     1,
		},
	}
	_, err := c.model.PermissionModel.AddMenus(ms)
	if err != nil {
		return err
	}
	return nil
}
