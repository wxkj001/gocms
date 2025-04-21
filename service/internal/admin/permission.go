package admin

import (
	"gocms/model"
	"gocms/router"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
)

type permission struct {
	*AdminRouter
}

func NewPermission(admin *AdminRouter) *permission {
	return &permission{AdminRouter: admin}
}
func (p *permission) List(ctx *gin.Context) {
	id := ctx.GetFloat64("role_id")
	isSuper := ctx.GetFloat64("is_super")
	if isSuper == 1 {
		id = -1
	}
	permissions, err := p.model.PermissionModel.GetPermissionsByRoleID(int64(id))
	if err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(200, router.Response{
		Code: 200,
		Data: permissions,
	})
}

// tree
func (p *permission) Tree(ctx *gin.Context) {
	id := ctx.GetFloat64("role_id")
	isSuper := int(ctx.GetFloat64("is_super"))
	if isSuper == 1 {
		id = -1
	}
	permissions, err := p.model.PermissionModel.GetPermissionsTreeAll(int64(id))
	if err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(200, router.Response{
		Code: 200,
		Data: permissions,
	})
}

// 根据type获取列表
func (p *permission) GetTypeList(ctx *gin.Context) {
	t := ctx.Param("type")
	id := ctx.Param("id")
	iid, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	menu, err := p.model.PermissionModel.GetMenuTreeByType(t, iid)
	if err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(200, router.Response{
		Code: 200,
		Data: menu,
	})
}

// 新增
func (p *permission) Add(ctx *gin.Context) {
	permission := &model.Permission{}
	err := ctx.ShouldBindBodyWith(permission, binding.JSON)
	if err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	_, err = p.model.PermissionModel.AddMenu(permission)
	if err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(200, router.Response{
		Code: 200,
		Data: permission,
	})

}

// 修改
func (p *permission) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	iid, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	permission := &model.Permission{}
	err = ctx.ShouldBindBodyWith(permission, binding.JSON)
	if err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	permission.ID = int64(iid)
	_, err = p.model.PermissionModel.UpdateMenu(permission)
	if err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(200, router.Response{
		Code: 200,
		Data: permission,
	})
}

// 删除
func (p *permission) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	iid, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	err = p.model.PermissionModel.DeleteMenu(int64(iid))
	if err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(200, router.Response{
		Code: 200,
		Data: iid,
	})
}
func (p *permission) Routers(ctx *gin.Context) {
	list := []string{}
	for _, v := range p.g.Routes() {
		p.log.Info(v.Path, zap.String("method", v.Method))
		list = append(list, v.Path)
	}
	ctx.JSON(200, router.Response{
		Code: 200,
		Data: list,
	})
}
