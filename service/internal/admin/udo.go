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
type udo struct {
	*AdminRouter
}

func NewUdo(admin *AdminRouter) *udo {
	return &udo{AdminRouter: admin}
}

func (c *udo) GetUdoObjectList(ctx *gin.Context) {
	list, err := c.model.UdoObjectModel.ObjectList(map[string]string{})
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: list,
	})
}
func (c *udo) GetUdoObject(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: "id不能为空",
		})
		return
	}
	iid, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	object, err := c.model.UdoObjectModel.GetUdoObjectByID(int64(iid))
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: object,
	})
}

// 创建
func (c *udo) CreateUdoObject(ctx *gin.Context) {
	object := &model.UdoObject{}
	err := ctx.ShouldBindJSON(object)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	err = c.model.UdoObjectModel.CreatedObject(object)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: object,
	})
}

// 更新
func (c *udo) UpdateUdoObject(ctx *gin.Context) {
	object := &model.UdoObject{}
	err := ctx.ShouldBindJSON(object)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	err = c.model.UdoObjectModel.UpdateObject(object)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: object,
	})
}
