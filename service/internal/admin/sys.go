package admin

import (
	"gocms/model"
	"gocms/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

type sysConfig struct {
	*AdminRouter
}

func NewSysConfig(admin *AdminRouter) *sysConfig {
	return &sysConfig{AdminRouter: admin}
}
func (s *sysConfig) List(ctx *gin.Context) {
	list, err := s.model.SysConfigModel.GetAllConfig()
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
func (s *sysConfig) Update(ctx *gin.Context) {
	var config model.SysConfig
	err := ctx.ShouldBindJSON(&config)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	err = s.model.SysConfigModel.UpdateConfig(&config)
	if err != nil {
		ctx.JSON(http.StatusOK, router.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, router.Response{
		Code: 200,
		Data: config,
	})
}
