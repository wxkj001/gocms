package admin

import (
	"gocms/model"
	"gocms/router"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type role struct {
	*AdminRouter
}

func NewRole(admin *AdminRouter) *role {
	return &role{AdminRouter: admin}
}

// tree
func (a *role) Tree(ctx *gin.Context) {
	list, err := a.model.RoleModel.GetListAll()
	if err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
			Data:    nil,
		})
	}
	ctx.JSON(200, router.Response{
		Code: 200,
		Data: list,
	})
}

// 列表
func (a *role) List(ctx *gin.Context) {
	pageNoReq := ctx.DefaultQuery("pageNo", "1")
	pageSizeReq := ctx.DefaultQuery("pageSize", "10")
	// enable := ctx.Query("enable")
	name := ctx.Query("name")
	pageNo, _ := strconv.Atoi(pageNoReq)
	pageSize, _ := strconv.Atoi(pageSizeReq)
	list, count, err := a.model.RoleModel.GetList(pageSize, pageNo, name)
	if err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	ctx.JSON(200, router.Response{
		Code: 200,
		Data: gin.H{
			"pageData": list,
			"total":    count,
		},
	})
}

// 添加
func (a *role) Add(ctx *gin.Context) {
	type RoleReq struct {
		*model.Role
		PermissionIds []int64 `json:"permissionIds"`
	}
	role := &RoleReq{}
	if err := ctx.ShouldBindBodyWith(role, binding.JSON); err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	_, err := a.model.RoleModel.CreateRole(role.Role)
	if err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	plist, err := a.model.PermissionModel.GetPermissionsByRoleIDAndType(role.PermissionIds, "API")
	if err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	for _, v := range plist {
		a.Casbin.AddPolicy(strconv.Itoa(int(role.Role.ID)), v.Path, "admin")
	}
	for _, v := range role.PermissionIds {
		a.Casbin.AddPolicy(strconv.Itoa(int(role.Role.ID)), strconv.Itoa(int(v)), "role")
	}
	a.Casbin.LoadPolicy()
	a.Casbin.SavePolicy()
	ctx.JSON(200, router.Response{
		Code: 200,
	})
}

// 删除
func (a *role) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: "id不能为空",
			Data:    nil,
		})
		return
	}
	iid, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	_, err = a.model.RoleModel.DeleteRole(int64(iid))
	if err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	err = a.model.RuleModel.Delete(&model.Rule{
		Ptype: "p",
		V0:    id,
		V1:    "",
		V2:    "admin",
		V3:    "",
		V4:    "",
		V5:    "",
	})
	if err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	err = a.model.RuleModel.Delete(&model.Rule{
		Ptype: "p",
		V0:    id,
		V1:    "",
		V2:    "role",
		V3:    "",
		V4:    "",
		V5:    "",
	})
	if err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	a.Casbin.LoadPolicy()
	a.Casbin.SavePolicy()
	ctx.JSON(200, router.Response{
		Code: 200,
		Data: nil,
	})
}

// 修改
func (a *role) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	iid, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	type RoleReq struct {
		*model.Role
		PermissionIds []int64 `json:"permissionIds"`
	}
	role := &RoleReq{}
	if err := ctx.ShouldBindBodyWith(role, binding.JSON); err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	role.ID = int64(iid)
	_, err = a.model.RoleModel.UpdateRole(role.Role)
	if err != nil {
		ctx.JSON(200, router.Response{
			Code:    500,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	if len(role.PermissionIds) > 1 {
		err = a.model.RuleModel.Delete(&model.Rule{
			Ptype: "p",
			V0:    id,
			V1:    "",
			V2:    "",
			V3:    "",
			V4:    "",
			V5:    "",
		})
		if err != nil {
			ctx.JSON(200, router.Response{
				Code:    500,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
		err = a.model.RuleModel.Delete(&model.Rule{
			Ptype: "p",
			V0:    id,
			V1:    "",
			V2:    "role",
			V3:    "",
			V4:    "",
			V5:    "",
		})
		if err != nil {
			ctx.JSON(200, router.Response{
				Code:    500,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
		a.Casbin.LoadPolicy()
		a.Casbin.SavePolicy()
		plist, err := a.model.PermissionModel.GetPermissionsByRoleIDAndType(role.PermissionIds, "API")
		if err != nil {
			ctx.JSON(200, router.Response{
				Code:    500,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
		for _, v := range plist {
			a.Casbin.AddPolicy(strconv.Itoa(int(role.Role.ID)), v.Path, v.Method)
		}
		for _, v := range role.PermissionIds {
			a.Casbin.AddPolicy(strconv.Itoa(int(role.Role.ID)), strconv.Itoa(int(v)), "role")
		}
		a.Casbin.LoadPolicy()
		a.Casbin.SavePolicy()

	}
	ctx.JSON(200, router.Response{
		Code: 200,
		Data: nil,
	})
}
