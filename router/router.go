package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Route interface {
	RouteRegister(*gin.Engine, *gin.RouterGroup)
}

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}

type Response struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}
