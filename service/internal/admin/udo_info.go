package admin

import (
	"gocms/router"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type udoInfo struct {
	*AdminRouter
}

func NewUdoInfo(admin *AdminRouter) *udoInfo {
	return &udoInfo{AdminRouter: admin}
}

func (u *udoInfo) List(ctx *gin.Context) {
	objcode := ctx.Param("code")
	// 获取参数
	pageNoReq := ctx.DefaultQuery("pageNo", "1")
	pageSizeReq := ctx.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(pageNoReq)
	pageSize, _ := strconv.Atoi(pageSizeReq)
	where := ctx.QueryMap("where")
	// 获取对象
	count, list, err := u.model.UdoDataModel.GetListByCode(objcode, pageNo, pageSize, where)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, router.Response{
		Code:    200,
		Message: "success",
		Data: gin.H{
			"pageData": list,
			"total":    count,
		},
	})
}

// 新增
func (u *udoInfo) Created(ctx *gin.Context) {
	objcode := ctx.Param("code")
	// 获取参数
	params := make(map[string]interface{})
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	// 获取对象
	_, err := u.model.UdoDataModel.Created(objcode, params)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
	})
}

// 更新
func (u *udoInfo) Update(ctx *gin.Context) {
	objcode := ctx.Param("code")
	// 获取参数
	params := make(map[string]interface{})
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	// 获取对象
	id := ctx.Param("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)
	_, err := u.model.UdoDataModel.Update(objcode, idInt, params)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
	})
}

// 删除
func (u *udoInfo) Delete(ctx *gin.Context) {
	objcode := ctx.Param("code")
	// 获取参数
	id := ctx.Param("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)
	// 获取对象
	_, err := u.model.UdoDataModel.Delete(objcode, idInt)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
	})
}
