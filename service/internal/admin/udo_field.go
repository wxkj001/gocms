package admin

import (
	"gocms/model"
	"gocms/router"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UDO
// 路由
type udoField struct {
	*AdminRouter
}

func NewUdoField(admin *AdminRouter) *udoField {
	return &udoField{AdminRouter: admin}
}

// 通过code获取字段列表
func (u *udoField) GetFieldListByCode(ctx *gin.Context) {
	// 获取参数
	objcode := ctx.Param("code")
	obj, err := u.model.UdoObjectModel.GetUdoObjectByCode(objcode)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	// 获取列表
	list, err := u.model.UdoFieldModel.GetFieldListByObjectId(obj.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
	}
	// 返回数据
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: list,
	})
}

// 列表
func (u *udoField) List(ctx *gin.Context) {
	// 获取参数
	pageNoReq := ctx.DefaultQuery("pageNo", "1")
	pageSizeReq := ctx.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(pageNoReq)
	pageSize, _ := strconv.Atoi(pageSizeReq)
	objectId, _ := strconv.Atoi(ctx.Query("object_id"))
	// 获取列表
	count, list, err := u.model.UdoFieldModel.GetFieldListPageByObjectId(int64(objectId), pageNo, pageSize, ctx.QueryMap("where"))
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	// 返回数据
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: gin.H{
			"pageData": list,
			"total":    count,
		},
	})
}

// 创建
func (u *udoField) Created(ctx *gin.Context) {
	// 获取参数
	var data model.UdoField
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	// 创建
	err = u.model.UdoFieldModel.CreatedField(data)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	// 返回数据
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: data,
	})
}

// 更新
func (u *udoField) Update(ctx *gin.Context) {
	// 获取参数
	var data model.UdoField
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	// 更新
	err = u.model.UdoFieldModel.UpdateField(data)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	// 返回数据
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: data,
	})
}

// 删除
func (u *udoField) Delete(ctx *gin.Context) {
	// 获取参数
	id, _ := strconv.Atoi(ctx.Param("id"))
	// 删除
	err := u.model.UdoFieldModel.DeleteField(int64(id))
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
	}
	// 返回数据
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: nil,
	})
}
